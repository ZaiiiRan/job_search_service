package interfaces

import (
	"context"

	"github.com/ZaiiiRan/job_search_service/auth-service/internal/domain/token"
)

type TokenCacheRepository interface {
	Get(ctx context.Context, token string) (*token.Token, error)
	Set(ctx context.Context, token *token.Token) error
	Del(ctx context.Context, token string) error
}
