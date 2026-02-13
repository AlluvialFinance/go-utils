package sql

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
)

// Create a temporary DB scoped to the test scenario and returns the config to connect to it
// temporary DB will be automically cleaned up after test
func CreateTempDB(t *testing.T, cfg *Config) (*Config, error) {
	t.Helper()
	conn, err := PGXConnect(t.Context(), cfg)
	if err != nil {
		return nil, err
	}

	name := sanitizeName(t.Name())
	err = createDB(t.Context(), conn, name)
	if err != nil {
		return nil, err
	}

	t.Cleanup(func() {
		_ = dropDB(t.Context(), conn, name)
		conn.Close(t.Context())
	})

	dbCfg := new(Config)
	*dbCfg = *cfg
	dbCfg.DBName = name

	return dbCfg, nil
}

func createDB(ctx context.Context, conn *pgx.Conn, dbName string) error {
	_, err := conn.Exec(ctx, fmt.Sprintf("CREATE DATABASE %v", dbName))
	if err != nil {
		return err
	}

	return nil
}

func dropDB(ctx context.Context, conn *pgx.Conn, dbName string) error {
	_, err := conn.Exec(ctx, fmt.Sprintf("DROP DATABASE %v", dbName))
	if err != nil {
		return err
	}

	return nil
}

func sanitizeName(name string) string {
	name = strings.ReplaceAll(name, "/", "_")
	name = strings.ToLower(name)
	return name
}
