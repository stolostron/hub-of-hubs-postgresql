package compliance

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

const DefaultLeafHubsNumber = 1000

func RunUpdate(ctx context.Context, dbConnectionPool *pgxpool.Pool, leafsNumber int) error {
	entry := time.Now()

	defer func() {
		now := time.Now()
		elapsed := now.Sub(entry)
		log.Printf("compliance RunUpdate of %d leaf hubs: elapsed %v\n", leafsNumber, elapsed)
	}()

	return nil
}
