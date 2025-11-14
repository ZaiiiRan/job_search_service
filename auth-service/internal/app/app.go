package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ZaiiiRan/job_search_service/auth-service/internal/config"
	usergrpcclient "github.com/ZaiiiRan/job_search_service/auth-service/internal/transport/client/grpc/user_client"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/transport/postgres"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/transport/redis"
	grpcserver "github.com/ZaiiiRan/job_search_service/auth-service/internal/transport/server/grpc"
	httpgateway "github.com/ZaiiiRan/job_search_service/auth-service/internal/transport/server/http"
	"github.com/ZaiiiRan/job_search_service/common/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type App struct {
	cfg *config.ServerConfig
	log *zap.SugaredLogger

	postgresClient *postgres.PostgresClient
	redisClient    *redis.RedisClient

	userGrpcClient *usergrpcclient.Client

	grpcServer  *grpcserver.Server
	httpGateway *httpgateway.Server
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
	if err := a.initRedisClient(ctx); err != nil {
		return err
	}

	if err := a.initGrpcServer(); err != nil {
		return err
	}
	if err := a.initHttpGateway(ctx); err != nil {
		return err
	}

	if err := a.initUserGrpcClient(ctx); err != nil {
		return err
	}

	a.startGrpcServer()
	a.startHttpGateway()

	a.log.Infow("app.started")
	return nil
}

func (a *App) Stop(ctx context.Context) {
	a.log.Infow("app.stopping")

	shCtx, cancel := context.WithTimeout(ctx, time.Duration(a.cfg.Shutdown.ShutdownTimeout)*time.Second)
	defer cancel()

	a.postgresClient.Close()
	a.redisClient.Close()
	a.grpcServer.Stop(shCtx)
	a.httpGateway.Stop(shCtx)

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

func (a *App) initUserGrpcClient(ctx context.Context) error {
	a.log.Infow("user grpc addr", "addr", a.cfg.UserServiceGRPCClient.Address)
	userClient, err := usergrpcclient.New(ctx, a.cfg.UserServiceGRPCClient, nil, nil)
	if err != nil {
		a.log.Errorw("app.user_grpc_client_init_failed", "err", err)
		return err
	}
	a.userGrpcClient = userClient

	if a.cfg.UserServiceGRPCClient.AutoConnect {
		a.log.Infow("app.user_grpc_client_connected")
	} else {
		a.log.Infow("app.user_grpc_client_initialized")
	}
	return nil
}

func (a *App) initGrpcServer() error {
	srv, err := grpcserver.New(a.cfg.GRPCServer, a.log)
	if err != nil {
		a.log.Errorw("app.grpc_server_init_failed", "err", err)
		return err
	}

	a.grpcServer = srv
	return nil
}

func (a *App) startGrpcServer() {
	go func() {
		a.log.Infow("app.grpc_serve_start", "port", a.cfg.GRPCServer.Port)
		if err := a.grpcServer.Start(); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			a.log.Fatalw("app.grpc_serve_error", "err", err)
		}
	}()
}

func (a *App) initHttpGateway(ctx context.Context) error {
	srv, err := httpgateway.New(ctx, a.cfg.HTTPGatewayServer, fmt.Sprintf("localhost%s", a.cfg.GRPCServer.Port))
	if err != nil {
		a.log.Errorw("app.http_gateway_init_failed", "err", err)
		return err
	}
	a.httpGateway = srv
	return nil
}

func (a *App) startHttpGateway() {
	go func() {
		a.log.Infow("app.http_gateway_start", "port", a.cfg.HTTPGatewayServer.Port)
		if err := a.httpGateway.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.log.Fatalw("app.http_gateway_error", "err", err)
		}
	}()
}
