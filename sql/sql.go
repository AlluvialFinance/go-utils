package sql

import (
	"context"
	"database/sql"
	"fmt"
)

func PingDB(ctx context.Context, db *sql.DB) error {
	err := db.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("health check failed on ping: %w", err)
	}

	rows, err := db.QueryContext(ctx, `SELECT VERSION()`)
	if err != nil {
		return fmt.Errorf("health check failed on select: %w", err)
	}

	if err = rows.Err(); err != nil {
		_ = rows.Close()
		return fmt.Errorf("health check failed after select: %w", err)
	}

	if err = rows.Close(); err != nil {
		return fmt.Errorf("health check failed on close: %w", err)
	}

	return nil
}
