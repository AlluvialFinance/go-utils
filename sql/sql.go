package sql

import (
	"context"
	"database/sql"
	"errors"
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

	if !rows.Next() {
		if err = rows.Err(); err != nil {
			_ = rows.Close()
			return fmt.Errorf("health check failed after select: %w", err)
		}
		_ = rows.Close()
		return errors.New("health check failed: no rows returned")
	}

	var version string
	if err = rows.Scan(&version); err != nil {
		_ = rows.Close()
		return fmt.Errorf("health check failed scanning version row: %w", err)
	}

	if err = rows.Err(); err != nil {
		_ = rows.Close()
		return fmt.Errorf("health check failed after scan: %w", err)
	}

	if err = rows.Close(); err != nil {
		return fmt.Errorf("health check failed on close: %w", err)
	}

	return nil
}
