package main

import (
	"fmt"
	"context"
	"os"
	"os/signal"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/open-cluster-management/hub-of-hubs-postgresql/performance_tests/pkg/compliance"
)

const (
	rowsNumber                           = "ROWS_NUMBER"
	environmentVariableDatabaseURL       = "DATABASE_URL"
)

func doMain() int {
	c := make(chan os.Signal, 1)
	signal.Notify(c)

	ctx, cancelContext := context.WithCancel(context.Background())
	defer cancelContext()

	go func() {
		<-c
		cancelContext()
	}()


	databaseURL, found := os.LookupEnv(environmentVariableDatabaseURL)
	if !found {
		fmt.Errorf("Not found environment variable %s\n", environmentVariableDatabaseURL)
		return 1
	}

	rowsNumber, found := os.LookupEnv(rowsNumber)
	if !found {
		fmt.Errorf("Not found environment variable %s\n", rowsNumber)
		return 1
	}

	dbConnectionPool, err := pgxpool.Connect(ctx, databaseURL)
	if err != nil {
		fmt.Errorf("Failed to connect to the database: %w", err)
		return 1
	}
	defer dbConnectionPool.Close()

	err = compliance.RunInsert(ctx, dbConnectionPool, 1000)
	if err != nil {
		fmt.Errorf("Failed to run compliance.RunInsert: %w", err)
		return 1
	}

	return 0
}

func main() {
	os.Exit(doMain())
}
