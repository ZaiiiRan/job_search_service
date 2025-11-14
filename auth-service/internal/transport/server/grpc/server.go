package grpcserver

import (
	"context"
	"fmt"
	"net"
	"time"

	pb "github.com/ZaiiiRan/job_search_service/auth-service/gen/go/auth_service/v1"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/config/settings"
	middleware "github.com/ZaiiiRan/job_search_service/common/pkg/middleware/grpc/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type Server struct {
	srv *grpc.Server
	lis net.Listener
}

func New(
	srvSettings settings.GRPCServerSettings,
	jwtSettings settings.JWTSettings,
	log *zap.SugaredLogger,
) (*Server, error) {
	s := grpc.NewServer(
		newChainUnaryInterceptor(log),
		grpc.KeepaliveParams(getGRPCKeepAliveServerParams(&srvSettings)),
		grpc.KeepaliveEnforcementPolicy(getGRPCKeepAliveEnforcement(&srvSettings)),
	)

	pb.RegisterAuthServiceServer(s, newAuthHandler())

	lis, err := net.Listen("tcp", srvSettings.Port)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	return &Server{
		srv: s,
		lis: lis,
	}, nil
}

func (s *Server) Start() error {
	return s.srv.Serve(s.lis)
}

func (s *Server) Stop(ctx context.Context) error {
	stopped := make(chan struct{})
	go func() {
		s.srv.GracefulStop()
		close(stopped)
	}()
	select {
	case <-ctx.Done():
		s.srv.Stop()
		return ctx.Err()
	case <-stopped:
		return nil
	}
}

func (s *Server) Addr() string {
	if s.lis != nil {
		return s.lis.Addr().String()
	}
	return ""
}

func newChainUnaryInterceptor(jwtSettings *settings.JWTSettings, log *zap.SugaredLogger) grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(
		middleware.RequestIdMiddleware(),
		middleware.LogMiddleware(log),
		middleware.RecoveryInterceptor(log),

		middleware.ApplicantAuthMiddleware(
			[]byte(jwtSettings.AccessTokenSecret),
			middleware.MiddlewareOnly(
				"/auth_service.v1.AuthService/GetNewApplicantActivationCode",
				"/auth_service.v1.AuthService/ActivateApplicant",
				"/auth_service.v1.AuthService/LogoutApplicant",
				"/auth_service.v1.AuthService/ChangeApplicantPassword",
			),
		),
		middleware.EmployerAuthMiddleware(
			[]byte(jwtSettings.AccessTokenSecret),
			middleware.MiddlewareOnly(
				"/auth_service.v1.AuthService/GetNewEmployerActivationCode",
				"/auth_service.v1.AuthService/ActivateEmployer",
				"/auth_service.v1.AuthService/LogoutEmployer",
				"/auth_service.v1.AuthService/ChangeEmployerPassword",
			),
		),
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
