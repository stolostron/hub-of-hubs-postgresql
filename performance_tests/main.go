package main

import (
	"fmt"
	"context"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/open-cluster-management/hub-of-hubs-postgresql/performance_tests/pkg/compliance"
)

const (
	environmentVariableRowsNumber                           = "ROWS_NUMBER"
	environmentVariableDatabaseURL       = "DATABASE_URL"
)

func doMain() int {
	ctx, cancelContext := context.WithCancel(context.Background())
	defer cancelContext()


	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		s := <-c
		fmt.Println("got signal", s)
		cancelContext()
	}()

	databaseURL, found := os.LookupEnv(environmentVariableDatabaseURL)
	if !found {
		fmt.Printf("Not found environment variable %s\n", environmentVariableDatabaseURL)
		return 1
	}

	rowsNumberString, found := os.LookupEnv(environmentVariableRowsNumber)
	if !found {
		fmt.Printf("Not found environment variable %s\n", environmentVariableRowsNumber)
		return 1
	}

	rowsNumber, err := strconv.Atoi(rowsNumberString)
	if err != nil {
		fmt.Printf("%s must be an integer\n", environmentVariableRowsNumber)
		return 1
	}

	if rowsNumber % 1000 != 0 {
		fmt.Printf("%s must be a multiple of 1000\n", environmentVariableRowsNumber)
		return 1
	}

	dbConnectionPool, err := pgxpool.Connect(ctx, databaseURL)
	if err != nil {
		fmt.Errorf("Failed to connect to the database: %w", err)
		return 1
	}
	defer dbConnectionPool.Close()

	err = compliance.RunInsert(ctx, dbConnectionPool, rowsNumber)
	if err != nil {
		fmt.Errorf("Failed to run compliance.RunInsert: %w", err)
		return 1
	}

	return 0
}

func main() {
	os.Exit(doMain())
}
