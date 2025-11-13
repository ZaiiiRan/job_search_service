package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ZaiiiRan/job_search_service/auth-service/internal/config/settings"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func Migrate(ctx context.Context, postgresCfg settings.PostgresSettings) error {
	db, err := sql.Open("pgx", postgresCfg.MigrationConnectionString)
	if err != nil {
		return fmt.Errorf("failed to open db connection: %w", err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping db: %w", err)
	}

	goose.SetTableName("goose_db_version")
	if err := goose.Up(db, "./migrations"); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	return nil
}
