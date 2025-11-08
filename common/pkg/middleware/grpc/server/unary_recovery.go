package middleware

import (
	"context"
	"runtime/debug"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func RecoveryInterceptor(log *zap.SugaredLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		defer func() {
			if r := recover(); r != nil {
				log.Errorw(
					"grpc.panic",
					"method", info.FullMethod,
					"panic", r,
					"stack", string(debug.Stack()),
				)
			}
		}()
		return handler(ctx, req)
	}
}
