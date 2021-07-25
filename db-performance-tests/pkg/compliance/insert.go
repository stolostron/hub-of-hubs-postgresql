package compliance

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	multipleInsertSize           = 1000
	copyInsertSize               = 10000
	columnSize                   = 7
	goRoutinesNumber             = 50
	maxNumberOfClusters          = 1000000
	maxNumberOfLeafHubs          = 1000
	compliantToNonCompliantRatio = 1000
)

func RunInsertByInsertWithMultipleValues(ctx context.Context, dbConnectionPool *pgxpool.Pool, n int) error {
	return doRunInsert(ctx, dbConnectionPool, n, insertRowsByInsertWithMultipleValues, "INSERT with multiple values", multipleInsertSize)
}

func RunInsertByCopy(ctx context.Context, dbConnectionPool *pgxpool.Pool, n int) error {
	return doRunInsert(ctx, dbConnectionPool, n, insertRowsByCopy, "COPY", copyInsertSize)
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
		fmt.Printf("compliance RunInsert %s: elapsed %v\n", description, elapsed)
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
			fmt.Printf("failed to insert rows: %v\n", err)
			break
		}
	}
}

func insertRowsByInsertWithMultipleValues(ctx context.Context, dbConnectionPool *pgxpool.Pool, insertSize int) error {
	rows := make([]interface{}, 0, insertSize*columnSize)

	for i := 0; i < insertSize; i++ {
		row, err := generateRow()
		if err != nil {
			return fmt.Errorf("failed to generate row: %w", err)
		}

		rows = append(rows, row...)
	}

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

	_, err := dbConnectionPool.Exec(ctx, sb.String(), rows...)
	if err != nil {
		return fmt.Errorf("insert into database failed: %w", err)
	}

	return nil
}

/* #nosec G404: Use of weak random number generator (math/rand instead of crypto/rand) */
func generateRow() ([]interface{}, error) {
	policyID, err := uuid.NewV4()
	if err != nil {
		return nil, fmt.Errorf("failed to generate UUID: %w", err)
	}

	clusterName := fmt.Sprintf("cluster%d", rand.Intn(maxNumberOfClusters))
	leafHubName := fmt.Sprintf("hub%d", rand.Intn(maxNumberOfLeafHubs))

	errorValue := "none"
	compliance := "compliant"

	if rand.Intn(compliantToNonCompliantRatio) == 0 {
		compliance = "non_compliant"
	}

	action := "inform"
	resourceVersion := strconv.Itoa(rand.Int())

	return []interface{}{policyID, clusterName, leafHubName, errorValue, compliance, action, resourceVersion}, nil
}

func insertRowsByCopy(ctx context.Context, dbConnectionPool *pgxpool.Pool, insertSize int) error {
	rows := make([][]interface{}, 0, insertSize)

	for i := 0; i < insertSize; i++ {
		row, err := generateRow()
		if err != nil {
			return fmt.Errorf("failed to generate row: %w", err)
		}

		rows = append(rows, row)
	}

	_, err := dbConnectionPool.CopyFrom(ctx, pgx.Identifier{"status", "compliance"},
		[]string{"policy_id", "cluster_name", "leaf_hub_name", "error", "compliance", "enforcement", "resource_version"},
		pgx.CopyFromRows(rows))
	if err != nil {
		return fmt.Errorf("insert into database failed: %w", err)
	}

	return nil
}
