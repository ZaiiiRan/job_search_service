package app

import (
	"context"
	"time"

	"github.com/ZaiiiRan/job_search_service/auth-service/internal/config"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/transport/postgres"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/transport/redis"
	"github.com/ZaiiiRan/job_search_service/common/pkg/logger"
	"go.uber.org/zap"
)

type App struct {
	cfg *config.ServerConfig
	log *zap.SugaredLogger

	postgresClient *postgres.PostgresClient
	redisClient    *redis.RedisClient
}

func New() (*App, error) {
	cfg, err := config.LoadServerConfig()
	if err != nil {
		return nil, err
	}

	log, err := logger.New()
	if err != nil {
		return nil, err
	}

	return &App{cfg: cfg, log: log}, nil
}

func (a *App) Run(ctx context.Context) error {
	if err := a.initPostgresClient(ctx); err != nil {
		return err
	}
	if err := a.initPostgresClient(ctx); err != nil {
		return err
	}

	a.log.Infow("app.started")
	return nil
}

func (a *App) Stop(ctx context.Context) {
	a.log.Infow("app.stopping")

	_, cancel := context.WithTimeout(ctx, time.Duration(a.cfg.Shutdown.ShutdownTimeout)*time.Second)
	defer cancel()

	a.postgresClient.Close()
	a.redisClient.Close()

	a.log.Infow("app.stopped")
}

func (a *App) initPostgresClient(ctx context.Context) error {
	if a.cfg.Migrate.NeedToMigrate {
		err := postgres.Migrate(ctx, a.cfg.DB)
		if err != nil {
			a.log.Errorw("app.postgres_migrate_failed", "err", err)
			return err
		}
	} else {
		a.log.Infow("app.postgres_migrate_skipped")
	}

	pgClient, err := postgres.New(ctx, a.cfg.DB)
	if err != nil {
		a.log.Errorw("app.postgres_connect_failed", "err", err)
		return err
	}
	a.postgresClient = pgClient

	a.log.Infow("app.postgres_connectd")
	return nil
}

func (a *App) initRedisClient(ctx context.Context) error {
	redisClient, err := redis.New(ctx, a.cfg.Redis)
	if err != nil {
		a.log.Errorw("app.redis_connect_failed", "err", err)
		return err
	}
	a.redisClient = redisClient

	a.log.Infow("app.redis_connected")
	return nil
}
