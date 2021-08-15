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
	clustersPerLeafHub           = 1000
	compliantToNonCompliantRatio = 1000
	DefaultRowsNumber            = 100000
	DefaultBatchSize             = 2000
	compliantString              = "compliant"
	nonCompliantString           = "non_compliant"
)

func RunInsertByInsertWithMultipleValues(ctx context.Context, dbConnectionPool *pgxpool.Pool, leafHubsNumber,
	batchSize int) error {
	return doRunInsert(ctx, dbConnectionPool, leafHubsNumber, insertRowsByInsertWithMultipleValues,
		"INSERT with multiple values", batchSize)
}

func RunInsertByCopy(ctx context.Context, dbConnectionPool *pgxpool.Pool, leafHubsNumber, batchSize int) error {
	return doRunInsert(ctx, dbConnectionPool, leafHubsNumber, insertRowsByCopy, "COPY", batchSize)
}

func doRunInsert(ctx context.Context, dbConnectionPool *pgxpool.Pool, leafHubsNumber int,
	insertFunc func(context.Context, *pgxpool.Pool, int, int) error, description string, batchSize int) error {
	_, err := dbConnectionPool.Exec(ctx, "DELETE from status.compliance")
	if err != nil {
		return fmt.Errorf("failed to clean the table before the test: %w", err)
	}

	entry := time.Now()

	defer func() {
		now := time.Now()
		elapsed := now.Sub(entry)
		log.Printf("compliance RunInsert %s %d leaf hubs by batch of %d rows: elapsed %v\n", description, leafHubsNumber,
			batchSize, elapsed)
	}()

	var wg sync.WaitGroup

	c := make(chan int, leafHubsNumber)

	for i := 0; i < leafHubsNumber; i++ {
		wg.Add(1)

		go insertRows(ctx, dbConnectionPool, c, &wg, insertFunc, batchSize)
	}

	for i := 0; i < leafHubsNumber; i++ {
		c <- i
	}
	close(c)

	wg.Wait()

	return nil
}

func insertRows(ctx context.Context, dbConnectionPool *pgxpool.Pool, c chan int, wg *sync.WaitGroup,
	insertFunc func(context.Context, *pgxpool.Pool, int, int) error, batchSize int) {
	defer wg.Done()

	for leafHubIndex := range c {
		if err := insertFunc(ctx, dbConnectionPool, leafHubIndex, batchSize); err != nil {
			log.Printf("failed to insert rows for leafHub %d: %v\n", leafHubIndex, err)
			break
		}
	}
}

func insertRowsByInsertWithMultipleValues(ctx context.Context, dbConnectionPool *pgxpool.Pool, leafHubIndex,
	batchSize int) error {
	rows := make([]interface{}, 0, clustersPerLeafHub*policiesNumber*columnSize)

	for clusterIndex := 0; clusterIndex < clustersPerLeafHub; clusterIndex++ {
		for policyIndex := 0; policyIndex < policiesNumber; policyIndex++ {
			rows = append(rows, generateRow(leafHubIndex, clusterIndex, policyIndex)...)
		}
	}

	var wg sync.WaitGroup

	batchesNumber := len(rows) / columnSize / batchSize
	c := make(chan []interface{}, batchesNumber)

	for i := 0; i < batchesNumber; i++ {
		wg.Add(1)

		go insertRowsBatchByInsertWithMultipleValues(ctx, dbConnectionPool, c, &wg)
	}

	for i := 0; i < batchesNumber; i++ {
		c <- rows[i*batchSize*columnSize : (i+1)*batchSize*columnSize]
	}
	close(c)

	wg.Wait()

	return nil
}

func insertRowsBatchByInsertWithMultipleValues(ctx context.Context, dbConnectionPool *pgxpool.Pool,
	c chan []interface{}, wg *sync.WaitGroup) {
	defer wg.Done()

	for rows := range c {
		err := doInsertRowsByInsertWithMultipleValues(ctx, dbConnectionPool, rows)
		if err != nil {
			log.Printf("failed to insert rows: %v\n", err)
		}
	}
}

func doInsertRowsByInsertWithMultipleValues(ctx context.Context, dbConnectionPool *pgxpool.Pool,
	rows []interface{}) error {
	sb := generateInsertByMultipleValues(len(rows) / columnSize)
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
func generateRow(leafHubIndex, clusterIndex, policyIndex int) []interface{} {
	policyID := policyUUIDs[policyIndex]
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

func insertRowsByCopy(ctx context.Context, dbConnectionPool *pgxpool.Pool, leafHubIndex, batchSize int) error {
	rows := make([][]interface{}, 0, clustersPerLeafHub*policiesNumber)

	for clusterIndex := 0; clusterIndex < clustersPerLeafHub; clusterIndex++ {
		for policyIndex := 0; policyIndex < policiesNumber; policyIndex++ {
			rows = append(rows, generateRow(leafHubIndex, clusterIndex, policyIndex))
		}
	}

	for i := 0; i < len(rows)/batchSize; i++ {
		_, err := dbConnectionPool.CopyFrom(ctx, pgx.Identifier{"status", "compliance"},
			[]string{"policy_id", "cluster_name", "leaf_hub_name", "error", "compliance", "enforcement", "resource_version"},
			pgx.CopyFromRows(rows[i*batchSize:(i+1)*batchSize]))
		if err != nil {
			return fmt.Errorf("insert into database failed: %w", err)
		}
	}

	return nil
}
