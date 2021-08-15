package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/open-cluster-management/hub-of-hubs-postgresql/performance_tests/pkg/compliance"
)

const (
	environmentVariableBatchSize            = "BATCH_SIZE"
	environmentVariableDatabaseURL          = "DATABASE_URL"
	environmentVariableInsertMultipleValues = "INSERT_MULTIPLE_VALUES"
	environmentVariableInsertCopy           = "INSERT_COPY"
	environmentVariableUpdate               = "UPDATE"
	environmentVariableLeafHubsNumber       = "LEAF_HUBS_NUMBER"
)

var (
	errEnvironmentVariableNotFound  = errors.New("not found environment variable")
	errEnvironmentVariableWrongType = errors.New("wrong type of environment variable")
)

func doMain() int {
	ctx, cancelContext := context.WithCancel(context.Background())
	defer cancelContext()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		s := <-c
		log.Println("got signal", s)
		cancelContext()
	}()

	databaseURL, batchSize, insertMultipleValues, insertCopy, update, leafHubsNumber,
		err := readEnvironmentVariables()
	if err != nil {
		log.Panicf("Failed to read environment variables: %v\n", err)
		return 1
	}

	dbConnectionPool, err := pgxpool.Connect(ctx, databaseURL)
	if err != nil {
		log.Printf("Failed to connect to the database: %v\n", err)
		return 1
	}
	defer dbConnectionPool.Close()

	err = runTests(ctx, dbConnectionPool, batchSize, leafHubsNumber, insertMultipleValues, insertCopy, update)
	if err != nil {
		log.Printf("Failed to run tests: %v\n", err)
		return 1
	}

	return 0
}

// long function of simple environment read, can be long.
func readEnvironmentVariables() (string, int, bool, bool, bool, int, error) {
	databaseURL, found := os.LookupEnv(environmentVariableDatabaseURL)

	if !found {
		return "", 0, false, false, false, 0,
			fmt.Errorf("%w: %s", errEnvironmentVariableNotFound, environmentVariableDatabaseURL)
	}

	batchSizeString, found := os.LookupEnv(environmentVariableBatchSize)
	if !found {
		batchSizeString = fmt.Sprintf("%d", compliance.DefaultBatchSize)
	}

	batchSize, err := strconv.Atoi(batchSizeString)
	if err != nil {
		return "", 0, false, false, false, 0,
			fmt.Errorf("%w: %s must be an integer", errEnvironmentVariableWrongType, environmentVariableBatchSize)
	}

	insertMultipleValues := false
	_, found = os.LookupEnv(environmentVariableInsertMultipleValues)

	if found {
		insertMultipleValues = true
	}

	insertCopy := false
	_, found = os.LookupEnv(environmentVariableInsertCopy)

	if found {
		insertCopy = true
	}

	update := false
	_, found = os.LookupEnv(environmentVariableUpdate)

	if found {
		update = true
	}

	leafHubsNumberString, found := os.LookupEnv(environmentVariableLeafHubsNumber)
	if !found {
		leafHubsNumberString = fmt.Sprintf("%d", compliance.DefaultLeafHubsNumber)
	}

	leafHubsNumber, err := strconv.Atoi(leafHubsNumberString)
	if err != nil {
		return "", 0, false, false, false, 0, fmt.Errorf("%w: %s must be an integer", errEnvironmentVariableWrongType,
			environmentVariableLeafHubsNumber)
	}

	return databaseURL, batchSize, insertMultipleValues, insertCopy, update, leafHubsNumber, nil
}

func runTests(ctx context.Context, dbConnectionPool *pgxpool.Pool, batchSize, leafHubsNumber int,
	insertMultipleValues, insertCopy, update bool) error {
	if insertMultipleValues {
		err := compliance.RunInsertByInsertWithMultipleValues(ctx, dbConnectionPool, leafHubsNumber, batchSize)
		if err != nil {
			return fmt.Errorf("failed to run compliance.RunInsert: %w", err)
		}
	}

	if insertCopy {
		err := compliance.RunInsertByCopy(ctx, dbConnectionPool, leafHubsNumber, batchSize)
		if err != nil {
			return fmt.Errorf("failed to run compliance.RunInsert: %w", err)
		}
	}

	if update {
		err := compliance.RunUpdate(ctx, dbConnectionPool, leafHubsNumber)
		if err != nil {
			return fmt.Errorf("failed to run compliance.RunInsert: %w", err)
		}
	}

	return nil
}

func main() {
	os.Exit(doMain())
}
