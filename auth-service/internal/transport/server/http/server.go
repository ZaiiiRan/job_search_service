package httpgateway

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	pb "github.com/ZaiiiRan/job_search_service/auth-service/gen/go/auth_service/v1"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/config/settings"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	srv *http.Server
}

func New(ctx context.Context, cfg settings.HTTPServerSettings, grpcAddr string) (*Server, error) {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := pb.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts); err != nil {
		return nil, fmt.Errorf("failed to register gateway handler: %w", err)
	}

	swaggerDir := filepath.Join("gen", "openapiv2", "auth_service", "v1")

	rootMux := http.NewServeMux()
	rootMux.Handle("/", mux)

	rootMux.Handle("/swagger/", http.StripPrefix("/swagger/",
		http.FileServer(http.Dir(swaggerDir)),
	))

	rootMux.Handle("/docs/",
		httpSwagger.Handler(
			httpSwagger.URL("/swagger/auth_service.swagger.json"),
			httpSwagger.DeepLinking(true),
			httpSwagger.DocExpansion("none"),
		),
	)

	srv := &http.Server{
		Addr:              cfg.Port,
		Handler:           rootMux,
		ReadHeaderTimeout: time.Duration(cfg.ReadHeaderTimeout) * time.Second,
		WriteTimeout:      time.Duration(cfg.WriteTimeout) * time.Second,
		ReadTimeout:       time.Duration(cfg.ReadTimeout) * time.Second,
		IdleTimeout:       time.Duration(cfg.IdleTimeout) * time.Second,
	}

	return &Server{srv: srv}, nil
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s *Server) Addr() string {
	return s.srv.Addr
}
