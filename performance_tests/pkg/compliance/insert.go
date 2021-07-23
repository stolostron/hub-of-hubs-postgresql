package compliance

import (
	"fmt"
	"context"
	"time"
	"github.com/jackc/pgx/v4/pgxpool"
)

func RunInsert(ctx context.Context, dbConnectionPool *pgxpool.Pool, n int) error {
	entry := time.Now()
	defer func() {
		now := time.Now()
		elapsed := now.Sub(entry)
		fmt.Printf("compliance RunInsert: elapsed %v\n", elapsed)
	}()

	_, err := dbConnectionPool.Exec(ctx, "DELETE from status.compliance")

	if err != nil {
		return fmt.Errorf("failed to clean the table before the test: %w", err)
	}

	return nil
}
