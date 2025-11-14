package middleware

import (
	"context"
	"errors"
	"strings"

	"github.com/ZaiiiRan/job_search_service/common/pkg/ctxmetadata"
	"github.com/ZaiiiRan/job_search_service/common/pkg/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errUnathorized = errors.New("unauthorized")
)

func ApplicantAuthMiddleware(secretKey []byte, shouldProtect MethodMatcher) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {

		if shouldProtect != nil && shouldProtect(info.FullMethod) {
			return handler(ctx, req)
		}

		tokenStr, err := extractBearerToken(ctx)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "%s", err.Error())
		}

		claims, err := jwt.ParseApplicantToken(tokenStr, secretKey)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "%s", errUnathorized.Error())
		}

		ctx = ctxmetadata.WithApplicantClaims(ctx, claims)

		return handler(ctx, req)
	}
}

func EmployerAuthMiddleware(secretKey []byte, shouldProtect MethodMatcher) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		if shouldProtect != nil && shouldProtect(info.FullMethod) {
			return handler(ctx, req)
		}

		tokenStr, err := extractBearerToken(ctx)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "%s", err.Error())
		}

		claims, err := jwt.ParseEmployerToken(tokenStr, secretKey)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "%s", errUnathorized.Error())
		}

		ctx = ctxmetadata.WithEmployerClaims(ctx, claims)

		return handler(ctx, req)
	}
}

func extractBearerToken(ctx context.Context) (string, error) {
	authHeader, err := ctxmetadata.GetAuthMetadataFromIncomingContext(ctx)
	if err != nil || len(authHeader) == 0 {
		return "", err
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", errUnathorized
	}

	return parts[1], nil
}
