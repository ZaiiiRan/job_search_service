package postgresimpl

import (
	"context"
	"fmt"
	"strings"

	userversion "github.com/ZaiiiRan/job_search_service/auth-service/internal/domain/user_version"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/impl"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/interfaces"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/models"
	postgresunitofwork "github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/unitofwork/postgres"
)

const (
	ApplicantVersion impl.RepositoryType = "applicant_version"
	EmployerVersion  impl.RepositoryType = "employer_version"
)

type UserVersionRepository struct {
	uow       *postgresunitofwork.UnitOfWork
	tableName string
}

func NewUserVersionRepository(uow *postgresunitofwork.UnitOfWork, repoType impl.RepositoryType) interfaces.UserVersionRepository {
	return &UserVersionRepository{uow: uow, tableName: string(repoType)}
}

func (r *UserVersionRepository) CreateUserVersion(ctx context.Context, uv *userversion.UserVersion) error {
	dal := models.V1UserVersionDalFromDomain(uv)

	conn, err := r.uow.GetConn(ctx)
	if err != nil {
		return err
	}

	var sb strings.Builder

	sb.WriteString(`
		INSERT INTO ` + r.tableName + ` (
			user_id,
			version,
			created_at,
			updated_at
		)
		SELECT
			(i).user_id,
			(i).version,
			(i).created_at,
			(i).updated_at
		FROM UNNEST($1::v1_user_version[]) i
		RETURNING
			id,
			user_id,
			version,
			created_at,
			updated_at
	`)

	rows, err := conn.Query(ctx, sb.String(), []models.V1UserVersionDal{dal})
	if err != nil {
		return fmt.Errorf("insert user version: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		var res models.V1UserVersionDal
		if err := rows.Scan(
			&res.Id,
			&res.UserId,
			&res.Version,
			&res.CreatedAt,
			&res.UpdatedAt,
		); err != nil {
			return fmt.Errorf("scan user version: %w", err)
		}
		*uv = *res.ToDomain()
		return nil
	}

	return fmt.Errorf("no user version returned from insert")
}

func (r *UserVersionRepository) UpdateUserVersion(ctx context.Context, uv *userversion.UserVersion) error {
	dal := models.V1UserVersionDalFromDomain(uv)

	conn, err := r.uow.GetConn(ctx)
	if err != nil {
		return err
	}

	var sb strings.Builder

	sb.WriteString(`
		UPDATE ` + r.tableName + ` AS t
		SET
			user_id = (i).user_id,
			version = (i).version,
			created_at = (i).created_at,
			updated_at = (i).updated_at
		FROM UNNEST($1::v1_user_version[]) AS i
		WHERE t.id = (i).id
		RETURNING
			t.id,
			t.user_id,
			t.version,
			t.created_at,
			t.updated_at
	`)

	rows, err := conn.Query(ctx, sb.String(), []models.V1UserVersionDal{dal})
	if err != nil {
		return fmt.Errorf("update user version: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		var res models.V1UserVersionDal
		if err := rows.Scan(
			&res.Id,
			&res.UserId,
			&res.Version,
			&res.CreatedAt,
			&res.UpdatedAt,
		); err != nil {
			return fmt.Errorf("scan user version: %w", err)
		}
		*uv = *res.ToDomain()
		return nil
	}

	return fmt.Errorf("no user version updated")
}

func (r *UserVersionRepository) DeleteUserVersion(ctx context.Context, uv *userversion.UserVersion) error {
	dal := models.V1UserVersionDalFromDomain(uv)

	conn, err := r.uow.GetConn(ctx)
	if err != nil {
		return err
	}

	var sb strings.Builder

	sb.WriteString(`
		DELETE FROM ` + r.tableName + ` AS t
		USING (
			SELECT (i).id AS id
			FROM UNNEST($1::v1_user_version[]) AS i
		) AS d
		WHERE t.id = d.id
		RETURNING
			t.id,
			t.user_id,
			t.version,
			t.created_at,
			t.updated_at
	`)

	rows, err := conn.Query(ctx, sb.String(), []models.V1UserVersionDal{dal})
	if err != nil {
		return fmt.Errorf("delete user version: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		var res models.V1UserVersionDal
		if err := rows.Scan(
			&res.Id,
			&res.UserId,
			&res.Version,
			&res.CreatedAt,
			&res.UpdatedAt,
		); err != nil {
			return fmt.Errorf("scan user version: %w", err)
		}
		*uv = *res.ToDomain()
		return nil
	}

	return fmt.Errorf("no user version deleted")
}

func (r *UserVersionRepository) QueryUserVersion(ctx context.Context, query *models.QueryUserVersionDal) (*userversion.UserVersion, error) {
	if query == nil {
		query = &models.QueryUserVersionDal{}
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
			id, user_id, version, created_at, updated_at
		FROM ` + r.tableName + `
		WHERE 1=1
	`)

	appendEqual(&sb, "id", query.Id, &args, &argPos)
	appendEqual(&sb, "user_id", query.UserId, &args, &argPos)
	appendLimitOffset(&sb, 1, 0, &args, &argPos)

	rows, err := conn.Query(ctx, sb.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("query user version: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		var res models.V1UserVersionDal
		if err := rows.Scan(
			&res.Id,
			&res.UserId,
			&res.Version,
			&res.CreatedAt,
			&res.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan user version: %w", err)
		}
		uv := res.ToDomain()
		return uv, nil
	}

	return nil, nil
}
