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
	environmentVariableRowsNumber           = "ROWS_NUMBER"
	environmentVariableBatchSize            = "BATCH_SIZE"
	environmentVariableDatabaseURL          = "DATABASE_URL"
	environmentVariableInsertMultipleValues = "INSERT_MULTIPLE_VALUES"
	environmentVariableInsertCopy           = "INSERT_COPY"
)

var (
	errEnvironmentVariableNotFound    = errors.New("not found environment variable")
	errEnvironmentVariableWrongType   = errors.New("wrong type of environment variable")
	errEnvironmentVariableWrongFormat = errors.New("wrong format of environment variable")
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

	databaseURL, rowsNumber, batchSize, insertMultipleValues, insertCopy, err := readEnvironmentVariables()
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

	err = runTests(ctx, dbConnectionPool, rowsNumber, batchSize, insertMultipleValues, insertCopy)
	if err != nil {
		log.Printf("Failed to run tests: %v\n", err)
		return 1
	}

	return 0
}

func readEnvironmentVariables() (string, int, int, bool, bool, error) {
	databaseURL, found := os.LookupEnv(environmentVariableDatabaseURL)

	if !found {
		return "", 0, 0, false, false, fmt.Errorf("%w: %s", errEnvironmentVariableNotFound, environmentVariableDatabaseURL)
	}

	rowsNumberString, found := os.LookupEnv(environmentVariableRowsNumber)
	if !found {
		return "", 0, 0, false, false, fmt.Errorf("%w: %s", errEnvironmentVariableNotFound, environmentVariableRowsNumber)
	}

	rowsNumber, err := strconv.Atoi(rowsNumberString)
	if err != nil {
		return "", 0, 0, false, false, fmt.Errorf("%w: %s must be an integer", errEnvironmentVariableWrongType,
			environmentVariableRowsNumber)
	}

	batchSizeString, found := os.LookupEnv(environmentVariableBatchSize)
	if !found {
		return "", 0, 0, false, false, fmt.Errorf("%w: %s", errEnvironmentVariableNotFound, environmentVariableBatchSize)
	}

	batchSize, err := strconv.Atoi(batchSizeString)
	if err != nil {
		return "", 0, 0, false, false, fmt.Errorf("%w: %s must be an integer", errEnvironmentVariableWrongType,
			environmentVariableBatchSize)
	}

	if rowsNumber%batchSize != 0 {
		return "", 0, 0, false, false, fmt.Errorf("%w: %s must be a multiple of %s", errEnvironmentVariableWrongFormat,
			environmentVariableRowsNumber, environmentVariableBatchSize)
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

	return databaseURL, rowsNumber, batchSize, insertMultipleValues, insertCopy, nil
}

func runTests(ctx context.Context, dbConnectionPool *pgxpool.Pool, rowsNumber, batchSize int, insertMultipleValues,
	insertCopy bool) error {
	if insertMultipleValues {
		err := compliance.RunInsertByInsertWithMultipleValues(ctx, dbConnectionPool, rowsNumber, batchSize)
		if err != nil {
			return fmt.Errorf("failed to run compliance.RunInsert: %w", err)
		}
	}

	if insertCopy {
		err := compliance.RunInsertByCopy(ctx, dbConnectionPool, rowsNumber, batchSize)
		if err != nil {
			return fmt.Errorf("failed to run compliance.RunInsert: %w", err)
		}
	}

	return nil
}

func main() {
	os.Exit(doMain())
}
