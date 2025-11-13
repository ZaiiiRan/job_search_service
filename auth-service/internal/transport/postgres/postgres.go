package postgres

import (
	"context"
	"fmt"

	"github.com/ZaiiiRan/job_search_service/auth-service/internal/config/settings"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresClient struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, postgresCfg settings.PostgresSettings) (*PostgresClient, error) {
	cfg, err := pgxpool.ParseConfig(postgresCfg.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("parse postgres config: %w", err)
	}

	cfg.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		names := []string{
			"v1_user_password", "_v1_user_password",
			"v1_activation_code", "_v1_activation_code",
			"v1_reset_password_code", "_v1_reset_password_code",
			"v1_refresh_token", "_v1_refresh_token",
		}
		types, err := conn.LoadTypes(ctx, names)
		if err != nil {
			return fmt.Errorf("load types: %w", err)
		}
		conn.TypeMap().RegisterTypes(types)
		return nil
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("new pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return &PostgresClient{pool: pool}, nil
}

func (p *PostgresClient) GetConn(ctx context.Context) (*pgxpool.Conn, error) {
	return p.pool.Acquire(ctx)
}

func (p *PostgresClient) Close() {
	if p.pool != nil {
		p.pool.Close()
	}
}
