package compliance

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

func RunUpdateAll(ctx context.Context, dbConnectionPool *pgxpool.Pool, leafsNumber, startLeafHubIndex int) error {
	entry := time.Now()

	defer func() {
		now := time.Now()
		elapsed := now.Sub(entry)
		log.Printf("compliance RunUpdateAll of %d leaf hubs: elapsed %v\n", leafsNumber, elapsed)
	}()

	var wg sync.WaitGroup

	c := make(chan int, leafsNumber)

	for i := 0; i < leafsNumber; i++ {
		wg.Add(1)

		go updateAllRows(ctx, dbConnectionPool, c, &wg)
	}

	for i := 0; i < leafsNumber; i++ {
		c <- startLeafHubIndex + i
	}
	close(c)

	wg.Wait()

	return nil
}

func updateAllRows(ctx context.Context, dbConnectionPool *pgxpool.Pool, c chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for leafHubIndex := range c {
		if err := updateAllRowsForLeafHub(ctx, dbConnectionPool, leafHubIndex); err != nil {
			log.Printf("failed to update rows: %v\n", err)
			break
		}
	}
}

func updateAllRowsForLeafHub(ctx context.Context, dbConnectionPool *pgxpool.Pool, leafHubIndex int) error {
	leafHubName := fmt.Sprintf("hub%d", leafHubIndex)

	_, err := dbConnectionPool.Exec(ctx,
		fmt.Sprintf("update status.compliance set compliance='compliant' where leaf_hub_name='%s'",
			leafHubName))
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	return nil
}
