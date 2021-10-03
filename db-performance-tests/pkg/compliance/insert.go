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
	columnSize                   = 6
	clustersPerLeafHub           = 1000
	compliantToNonCompliantRatio = 1000
	DefaultRowsNumber            = 100000
	DefaultBatchSize             = 2000
	DefaultLeafHubsNumber        = 1000
	DefaultStartLeafHubIndex     = 0
	compliantString              = "compliant"
	nonCompliantString           = "non_compliant"
	maxNumberOfGoRoutines        = 100
)

func RunInsertByInsertWithMultipleValues(ctx context.Context, dbConnectionPool *pgxpool.Pool, leafHubsNumber,
	startLeafHubIndex, batchSize int) error {
	return doRunInsert(ctx, dbConnectionPool, leafHubsNumber, startLeafHubIndex, insertRowsByInsertWithMultipleValues,
		"INSERT with multiple values", batchSize)
}

func RunInsertByCopy(ctx context.Context, dbConnectionPool *pgxpool.Pool, leafHubsNumber, startLeafHubIndex,
	batchSize int) error {
	return doRunInsert(ctx, dbConnectionPool, leafHubsNumber, startLeafHubIndex, insertRowsByCopy, "COPY", batchSize)
}

func doRunInsert(ctx context.Context, dbConnectionPool *pgxpool.Pool, leafHubsNumber, startLeafHubIndex int,
	insertFunc func(context.Context, *pgxpool.Pool, [][]interface{}) error, description string, batchSize int) error {
	entry := time.Now()

	defer func() {
		now := time.Now()
		elapsed := now.Sub(entry)
		log.Printf("compliance RunInsert %s %d leaf hubs by batch of %d rows: elapsed %v\n", description, leafHubsNumber,
			batchSize, elapsed)
	}()

	var wg sync.WaitGroup

	c := make(chan int, maxNumberOfGoRoutines)

	for i := 0; i < maxNumberOfGoRoutines; i++ {
		wg.Add(1)

		go insertRows(ctx, dbConnectionPool, c, &wg, insertFunc, batchSize)
	}

	for i := 0; i < leafHubsNumber; i++ {
		c <- startLeafHubIndex + i
	}
	close(c)

	wg.Wait()

	return nil
}

func insertRows(ctx context.Context, dbConnectionPool *pgxpool.Pool, c chan int, wg *sync.WaitGroup,
	insertFunc func(context.Context, *pgxpool.Pool, [][]interface{}) error, batchSize int) {
	defer wg.Done()

	for leafHubIndex := range c {
		insertRowsForLeafHub(ctx, dbConnectionPool, leafHubIndex, batchSize, insertFunc)
	}
}

func generateRowsForLeafHub(leafHubIndex int) [][]interface{} {
	rows := make([][]interface{}, 0, clustersPerLeafHub*policiesNumber)

	beginClusterIndex, endClusterIndex := leafHubIndex*clustersPerLeafHub, (leafHubIndex+1)*clustersPerLeafHub
	for clusterIndex := beginClusterIndex; clusterIndex < endClusterIndex; clusterIndex++ {
		for policyIndex := 0; policyIndex < policiesNumber; policyIndex++ {
			rows = append(rows, generateRow(leafHubIndex, clusterIndex, policyIndex))
		}
	}

	return rows
}

func insertRowsForLeafHub(ctx context.Context, dbConnectionPool *pgxpool.Pool, leafHubIndex,
	batchSize int, insertFunc func(context.Context, *pgxpool.Pool, [][]interface{}) error) {
	rows := generateRowsForLeafHub(leafHubIndex)

	var wg sync.WaitGroup

	batchesNumber := len(rows) / batchSize
	c := make(chan [][]interface{}, maxNumberOfGoRoutines)

	for i := 0; i < maxNumberOfGoRoutines; i++ {
		wg.Add(1)

		go insertRowsBatch(ctx, dbConnectionPool, c, &wg, insertFunc)
	}

	for i := 0; i < batchesNumber; i++ {
		c <- rows[i*batchSize : (i+1)*batchSize]
	}
	close(c)

	wg.Wait()
}

func insertRowsBatch(ctx context.Context, dbConnectionPool *pgxpool.Pool, c chan [][]interface{}, wg *sync.WaitGroup,
	insertFunc func(context.Context, *pgxpool.Pool, [][]interface{}) error) {
	defer wg.Done()

	for rows := range c {
		err := insertFunc(ctx, dbConnectionPool, rows)
		if err != nil {
			log.Printf("failed to insert rows: %v\n", err)
		}
	}
}

func flatten(rows [][]interface{}) []interface{} {
	resultRows := make([]interface{}, 0, len(rows)*columnSize)

	for _, row := range rows {
		resultRows = append(resultRows, row...)
	}

	return resultRows
}

func insertRowsByInsertWithMultipleValues(ctx context.Context, dbConnectionPool *pgxpool.Pool,
	rows [][]interface{}) error {
	sb := generateInsertByMultipleValues(len(rows))
	sb.WriteString(" ON CONFLICT DO NOTHING")

	_, err := dbConnectionPool.Exec(ctx, sb.String(), flatten(rows)...)
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

func generateRow(leafHubIndex, clusterIndex, policyIndex int) []interface{} {
	policyID := policyUUIDs[policyIndex]
	leafHubName := fmt.Sprintf("hub%d", leafHubIndex)
	clusterName := fmt.Sprintf("cluster%d", clusterIndex)

	errorValue, compliance, action := generateDerivedColumns()

	return []interface{}{policyID, clusterName, leafHubName, errorValue, compliance, action}
}

/* #nosec G404: Use of weak random number generator (math/rand instead of crypto/rand) */
func generateDerivedColumns() (string, string, string) {
	errorValue := "none"
	compliance := compliantString

	if rand.Intn(compliantToNonCompliantRatio) == 0 {
		compliance = nonCompliantString
	}

	action := "inform"

	return errorValue, compliance, action
}

func insertRowsByCopy(ctx context.Context, dbConnectionPool *pgxpool.Pool, rows [][]interface{}) error {
	_, err := dbConnectionPool.CopyFrom(ctx, pgx.Identifier{"status", "compliance"},
		[]string{"policy_id", "cluster_name", "leaf_hub_name", "error", "compliance", "enforcement"},
		pgx.CopyFromRows(rows))
	if err != nil {
		return fmt.Errorf("insert into database failed: %w", err)
	}

	return nil
}
