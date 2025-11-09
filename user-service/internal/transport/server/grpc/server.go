package grpcserver

import (
	"net"
	"time"

	middleware "github.com/ZaiiiRan/job_search_service/common/pkg/middleware/grpc/server"
	pb "github.com/ZaiiiRan/job_search_service/user-service/gen/go/user_service/v1"
	"github.com/ZaiiiRan/job_search_service/user-service/internal/config/settings"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func New(
	srvSettings settings.GRPCServerSettings,
	log *zap.SugaredLogger,
) (*grpc.Server, net.Listener, error) {
	s := grpc.NewServer(
		newChainUnaryInterceptor(log),
		grpc.KeepaliveParams(getGRPCKeepAliveServerParams(&srvSettings)),
		grpc.KeepaliveEnforcementPolicy(getGRPCKeepAliveEnforcement(&srvSettings)),
	)

	pb.RegisterUserServiceServer(s, newUserHandler())

	lis, err := net.Listen("tcp", srvSettings.Port)
	if err != nil {
		return nil, nil, err
	}

	return s, lis, err
}

func newChainUnaryInterceptor(log *zap.SugaredLogger) grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(
		middleware.RequestIdMiddleware(),
		middleware.LogMiddleware(log),
		middleware.RecoveryInterceptor(log),
	)
}

func getGRPCKeepAliveServerParams(c *settings.GRPCServerSettings) keepalive.ServerParameters {
	if c == nil {
		return keepalive.ServerParameters{}
	}
	return keepalive.ServerParameters{
		MaxConnectionIdle:     time.Duration(c.MaxConnectionIdle) * time.Second,
		MaxConnectionAge:      time.Duration(c.MaxConnectionAge) * time.Second,
		MaxConnectionAgeGrace: time.Duration(c.MaxConnectionAgeGrace) * time.Second,
		Time:                  time.Duration(c.KeepaliveTime) * time.Second,
		Timeout:               time.Duration(c.KeepaliveTimeout) * time.Second,
	}
}

func getGRPCKeepAliveEnforcement(c *settings.GRPCServerSettings) keepalive.EnforcementPolicy {
	if c == nil {
		return keepalive.EnforcementPolicy{}
	}
	return keepalive.EnforcementPolicy{
		MinTime:             0,
		PermitWithoutStream: c.PermitWithoutStream,
	}
}
