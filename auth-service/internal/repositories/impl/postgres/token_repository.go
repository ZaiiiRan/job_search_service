package postgresimpl

import (
	"context"
	"fmt"
	"strings"

	"github.com/ZaiiiRan/job_search_service/auth-service/internal/domain/token"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/impl"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/interfaces"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/models"
	postgresunitofwork "github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/unitofwork/postgres"
)

const (
	ApplicantRefreshTokenRepository impl.RepositoryType = "applicant_refresh_tokens"
	EmployerRefreshTokenRepository  impl.RepositoryType = "employer_refresh_tokens"
)

type TokenRepository struct {
	uow       *postgresunitofwork.UnitOfWork
	tableName string
}

func NewTokenRepository(uow *postgresunitofwork.UnitOfWork, repoType impl.RepositoryType) interfaces.TokenRepository {
	return &TokenRepository{uow: uow, tableName: string(repoType)}
}

func (r *TokenRepository) CreateToken(ctx context.Context, token *token.Token) error {
	dal := models.V1RefreshTokenDalFromDomain(token)

	conn, err := r.uow.GetConn(ctx)
	if err != nil {
		return err
	}

	var sb strings.Builder

	sb.WriteString(`
		INSERT INTO ` + r.tableName + ` (
			user_id,
			token,
			version,
			expires_at,
			created_at,
			updated_at
		)
		SELECT
			(i).user_id,
			(i).token,
			(i).version,
			(i).expires_at,
			(i).created_at,
			(i).updated_at
		FROM UNNEST($1::v1_refresh_token[]) i
		RETURNING
			id,
			user_id,
			token,
			version,
			expires_at,
			created_at,
			updated_at
	`)

	rows, err := conn.Query(ctx, sb.String(), []models.V1RefreshTokenDal{dal})
	if err != nil {
		return fmt.Errorf("insert token: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		var res models.V1RefreshTokenDal
		if err := rows.Scan(
			&res.Id,
			&res.UserId,
			&res.Token,
			&res.Version,
			&res.ExpiresAt,
			&res.CreatedAt,
			&res.UpdatedAt,
		); err != nil {
			return fmt.Errorf("scan token: %w", err)
		}
		*token = *res.ToDomain()
		return nil
	}

	return fmt.Errorf("no token returned from insert")
}

func (r *TokenRepository) UpdateToken(ctx context.Context, token *token.Token) error {
	dal := models.V1RefreshTokenDalFromDomain(token)

	conn, err := r.uow.GetConn(ctx)
	if err != nil {
		return err
	}

	var sb strings.Builder

	sb.WriteString(`
		UPDATE ` + r.tableName + ` AS t
		SET
			user_id = (i).user_id,
			token = (i).token,
			version = (i).version,
			expires_at = (i).expires_at,
			created_at = (i).created_at,
			updated_at = (i).updated_at
		FROM UNNEST($1::v1_refresh_token[]) i
		WHERE t.id = (i).id
		RETURNING
			t.id,
			t.user_id,
			t.token,
			t.version,
			t.expires_at,
			t.created_at,
			t.updated_at
	`)

	rows, err := conn.Query(ctx, sb.String(), []models.V1RefreshTokenDal{dal})
	if err != nil {
		return fmt.Errorf("update token: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		var res models.V1RefreshTokenDal
		if err := rows.Scan(
			&res.Id,
			&res.UserId,
			&res.Token,
			&res.Version,
			&res.ExpiresAt,
			&res.CreatedAt,
			&res.UpdatedAt,
		); err != nil {
			return fmt.Errorf("scan token: %w", err)
		}
		*token = *res.ToDomain()
		return nil
	}

	return fmt.Errorf("no token updated")
}

func (r *TokenRepository) DeleteToken(ctx context.Context, tokenStr string) error {
	conn, err := r.uow.GetConn(ctx)
	if err != nil {
		return err
	}

	var sb strings.Builder

	sb.WriteString(`
		DELETE FROM ` + r.tableName + ` AS t
		WHERE token = $1
	`)

	if _, err := conn.Exec(ctx, sb.String(), tokenStr); err != nil {
		return fmt.Errorf("delete token: %w", err)
	}
	return nil
}

func (r *TokenRepository) DeleteTokensByUserId(ctx context.Context, userId int64) error {
	conn, err := r.uow.GetConn(ctx)
	if err != nil {
		return err
	}

	var sb strings.Builder

	sb.WriteString(`
		DELETE FROM ` + r.tableName + ` AS t
		WHERE t.user_id = $1
	`)

	if _, err := conn.Exec(ctx, sb.String(), userId); err != nil {
		return fmt.Errorf("delete tokens by user id: %w", err)
	}
	return nil
}

func (r *TokenRepository) QueryToken(ctx context.Context, query *models.QueryTokenDal) (*token.Token, error) {
	if query == nil {
		query = &models.QueryTokenDal{}
	}

	conn, err := r.uow.GetConn(ctx)
	if err != nil {
		return nil, err
	}

	var (
		sb     strings.Builder
		args   []any
		argPos = 1
	)

	sb.WriteString(`
		SELECT
			id, user_id, token, version,
			expires_at, created_at, updated_at
		FROM ` + r.tableName + `
		WHERE 1=1
	`)

	appendEqual(&sb, "id", query.Id, &args, &argPos)
	appendEqual(&sb, "user_id", query.UserId, &args, &argPos)
	appendEqual(&sb, "token", query.Token, &args, &argPos)
	appendEqual(&sb, "version", query.Version, &args, &argPos)
	appendLimitOffset(&sb, 1, 0, &args, &argPos)

	rows, err := conn.Query(ctx, sb.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("query token: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		var res models.V1RefreshTokenDal
		if err := rows.Scan(
			&res.Id,
			&res.UserId,
			&res.Token,
			&res.Version,
			&res.ExpiresAt,
			&res.CreatedAt,
			&res.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan token: %w", err)
		}
		token := res.ToDomain()
		return token, nil
	}

	return nil, nil
}
