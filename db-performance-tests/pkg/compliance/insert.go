package compliance

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	columnSize                   = 7
	goRoutinesNumber             = 500
	clustersPerLeafHub           = 1000
	maxNumberOfLeafHubs          = 1000
	compliantToNonCompliantRatio = 1000
	DefaultRowsNumber            = 100000
	DefaultBatchSize             = 2000
	compliantString              = "compliant"
	nonCompliantString           = "non_compliant"
)

func RunInsertByInsertWithMultipleValues(ctx context.Context, dbConnectionPool *pgxpool.Pool, n, batchSize int) error {
	return doRunInsert(ctx, dbConnectionPool, n, insertRowsByInsertWithMultipleValues, "INSERT with multiple values",
		batchSize)
}

func RunInsertByCopy(ctx context.Context, dbConnectionPool *pgxpool.Pool, n, batchSize int) error {
	return doRunInsert(ctx, dbConnectionPool, n, insertRowsByCopy, "COPY", batchSize)
}

func doRunInsert(ctx context.Context, dbConnectionPool *pgxpool.Pool, n int,
	insertFunc func(context.Context, *pgxpool.Pool, int) error, description string, insertSize int) error {
	_, err := dbConnectionPool.Exec(ctx, "DELETE from status.compliance")
	if err != nil {
		return fmt.Errorf("failed to clean the table before the test: %w", err)
	}

	entry := time.Now()

	defer func() {
		now := time.Now()
		elapsed := now.Sub(entry)
		log.Printf("compliance RunInsert %s %d rows by batch of %d rows: elapsed %v\n", description, n, insertSize, elapsed)
	}()

	var wg sync.WaitGroup

	insertNumber := n / insertSize
	c := make(chan int, insertNumber)

	for i := 0; i < goRoutinesNumber; i++ {
		wg.Add(1)

		go insertRows(ctx, dbConnectionPool, c, &wg, insertFunc, insertSize)
	}

	for i := 0; i < insertNumber; i++ {
		c <- i
	}
	close(c)

	wg.Wait()

	return nil
}

func insertRows(ctx context.Context, dbConnectionPool *pgxpool.Pool, c chan int, wg *sync.WaitGroup,
	insertFunc func(context.Context, *pgxpool.Pool, int) error, insertSize int) {
	defer wg.Done()

	for range c {
		if err := insertFunc(ctx, dbConnectionPool, insertSize); err != nil {
			log.Printf("failed to insert rows: %v\n", err)
			break
		}
	}
}

func insertRowsByInsertWithMultipleValues(ctx context.Context, dbConnectionPool *pgxpool.Pool, insertSize int) error {
	rows := make([]interface{}, 0, insertSize*columnSize)

	for i := 0; i < insertSize; i++ {
		rows = append(rows, generateRow()...)
	}

	err := doInsertRowsByInsertWithMultipleValues(ctx, dbConnectionPool, rows, insertSize)
	if err != nil {
		return fmt.Errorf("insert into database failed: %w", err)
	}

	return nil
}

func doInsertRowsByInsertWithMultipleValues(ctx context.Context, dbConnectionPool *pgxpool.Pool, rows []interface{},
	insertSize int) error {
	sb := generateInsertByMultipleValues(insertSize)
	sb.WriteString(" ON CONFLICT DO NOTHING")

	_, err := dbConnectionPool.Exec(ctx, sb.String(), rows...)
	if err != nil {
		return fmt.Errorf("insert into database failed: %w", err)
	}

	return nil
}

func generateInsertByMultipleValues(insertSize int) *strings.Builder {
	var sb strings.Builder

	sb.WriteString("INSERT INTO status.compliance values")

	for i := 0; i < insertSize; i++ {
		sb.WriteString("(")

		for j := 0; j < columnSize; j++ {
			sb.WriteString("$")
			sb.WriteString(strconv.Itoa(i*columnSize + j + 1))

			if j < columnSize-1 {
				sb.WriteString(", ")
			}
		}

		sb.WriteString(")")

		if i < insertSize-1 {
			sb.WriteString(", ")
		}
	}

	return &sb
}

/* #nosec G404: Use of weak random number generator (math/rand instead of crypto/rand) */
func generateRow() []interface{} {
	policyIndex := rand.Intn(maxNumberOfPolicies)
	policyID := policyUUIDs[policyIndex]
	leafHubIndex := rand.Intn(maxNumberOfLeafHubs)
	clusterIndex := leafHubIndex*clustersPerLeafHub + rand.Intn(clustersPerLeafHub)
	leafHubName := fmt.Sprintf("hub%d", leafHubIndex)
	clusterName := fmt.Sprintf("cluster%d", clusterIndex)

	errorValue, compliance, action, resourceVersion := generateDerivedColumns(policyID.String(), leafHubName, clusterName)

	return []interface{}{policyID, clusterName, leafHubName, errorValue, compliance, action, resourceVersion}
}

/* #nosec G404: Use of weak random number generator (math/rand instead of crypto/rand) */
func generateDerivedColumns(policyID, leafHubName, clusterName string) (string, string, string, string) {
	errorValue := "none"
	compliance := compliantString

	if rand.Intn(compliantToNonCompliantRatio) == 0 {
		compliance = nonCompliantString
	}

	action := "inform"
	resourceVersion := fmt.Sprintf("%s%s%s", policyID, leafHubName, clusterName)

	return errorValue, compliance, action, resourceVersion
}

func insertRowsByCopy(ctx context.Context, dbConnectionPool *pgxpool.Pool, insertSize int) error {
	rows := make([][]interface{}, 0, insertSize)

	for i := 0; i < insertSize; i++ {
		rows = append(rows, generateRow())
	}

	_, err := dbConnectionPool.CopyFrom(ctx, pgx.Identifier{"status", "compliance"},
		[]string{"policy_id", "cluster_name", "leaf_hub_name", "error", "compliance", "enforcement", "resource_version"},
		pgx.CopyFromRows(rows))
	if err != nil {
		return fmt.Errorf("insert into database failed: %w", err)
	}

	return nil
}
