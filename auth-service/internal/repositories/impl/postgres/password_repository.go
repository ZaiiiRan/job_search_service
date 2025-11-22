package postgresimpl

import (
	"context"
	"fmt"
	"strings"

	"github.com/ZaiiiRan/job_search_service/auth-service/internal/domain/password"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/impl"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/interfaces"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/models"
	postgresunitofwork "github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/unitofwork/postgres"
)

const (
	ApplicantPasswordRepository impl.RepositoryType = "applicant_passwords"
	EmployerPasswordRepository  impl.RepositoryType = "employer_passwords"
)

type PasswordRepository struct {
	uow       *postgresunitofwork.UnitOfWork
	tableName string
}

func NewPasswordRepository(uow *postgresunitofwork.UnitOfWork, repoType impl.RepositoryType) interfaces.PasswordRepository {
	return &PasswordRepository{uow: uow, tableName: string(repoType)}
}

func (r *PasswordRepository) CreatePassword(ctx context.Context, password *password.Password) error {
	dal := models.V1UserPasswordDalFromDomain(password)

	conn, err := r.uow.GetConn(ctx)
	if err != nil {
		return err
	}

	var sb strings.Builder

	sb.WriteString(`
		INSERT INTO ` + r.tableName + ` (
			user_id,
			password_hash,
			created_at,
			updated_at
		) 
		SELECT 
			(i).user_id,
			(i).password_hash,
			(i).created_at,
			(i).updated_at
		FROM UNNEST($1::v1_user_password[]) i
		RETURNING
			id,
			user_id,
			password_hash,
			created_at,
			updated_at
	`)

	rows, err := conn.Query(ctx, sb.String(), []models.V1UserPasswordDal{dal})
	if err != nil {
		return fmt.Errorf("insert password: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		var res models.V1UserPasswordDal
		if err := rows.Scan(
			&res.Id,
			&res.UserId,
			&res.PassordHash,
			&res.CreatedAt,
			&res.UpdatedAt,
		); err != nil {
			return fmt.Errorf("scan password: %w", err)
		}
		*password = *res.ToDomain()
		return nil
	}

	return fmt.Errorf("no password returned from insert")
}

func (r *PasswordRepository) UpdatePassword(ctx context.Context, password *password.Password) error {
	dal := models.V1UserPasswordDalFromDomain(password)

	conn, err := r.uow.GetConn(ctx)
	if err != nil {
		return err
	}

	var sb strings.Builder

	sb.WriteString(`
		UPDATE ` + r.tableName + ` AS t
		SET
			user_id = (i).user_id,
			password_hash = (i).password_hash,
			created_at = (i).created_at,
			updated_at = (i).updated_at
		FROM UNNEST($1::v1_user_password[]) i
		WHERE t.id = (i).id
		RETURNING
			t.id,
			t.user_id,
			t.password_hash,
			t.created_at,
			t.updated_at
	`)

	rows, err := conn.Query(ctx, sb.String(), []models.V1UserPasswordDal{dal})
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		var res models.V1UserPasswordDal
		if err := rows.Scan(
			&res.Id,
			&res.UserId,
			&res.PassordHash,
			&res.CreatedAt,
			&res.UpdatedAt,
		); err != nil {
			return fmt.Errorf("scan password: %w", err)
		}
		*password = *res.ToDomain()
		return nil
	}

	return fmt.Errorf("no password updated")
}

func (r *PasswordRepository) DeletePassword(ctx context.Context, password *password.Password) error {
	dal := models.V1UserPasswordDalFromDomain(password)

	conn, err := r.uow.GetConn(ctx)
	if err != nil {
		return err
	}

	var sb strings.Builder

	sb.WriteString(`
		DELETE FROM ` + r.tableName + ` AS t
		USING (
			SELECT (i).id AS id
			FROM UNNEST($1::v1_user_password[]) AS i
		) AS d
		WHERE t.id = d.id
		RETURNING
			t.id,
			t.user_id,
			t.password_hash,
			t.created_at,
			t.updated_at
	`)

	rows, err := conn.Query(ctx, sb.String(), []models.V1UserPasswordDal{dal})
	if err != nil {
		return fmt.Errorf("delete password: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		var res models.V1UserPasswordDal
		if err := rows.Scan(
			&res.Id,
			&res.UserId,
			&res.PassordHash,
			&res.CreatedAt,
			&res.UpdatedAt,
		); err != nil {
			return fmt.Errorf("scan password: %w", err)
		}
		*password = *res.ToDomain()
		return nil
	}

	return fmt.Errorf("no password deleted")
}

func (r *PasswordRepository) QueryPassword(ctx context.Context, query *models.QueryPasswordDal) (*password.Password, error) {
	if query == nil {
		query = &models.QueryPasswordDal{}
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
			id, user_id, password_hash,
			created_at, updated_at
		FROM ` + r.tableName + `
		WHERE 1=1
	`)

	appendEqual(&sb, "id", query.Id, &args, &argPos)
	appendEqual(&sb, "user_id", query.UserId, &args, &argPos)
	appendLimitOffset(&sb, 1, 0, &args, &argPos)

	rows, err := conn.Query(ctx, sb.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("query password: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		var res models.V1UserPasswordDal
		if err := rows.Scan(
			&res.Id,
			&res.UserId,
			&res.PassordHash,
			&res.CreatedAt,
			&res.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan password: %w", err)
		}
		password := res.ToDomain()
		return password, nil
	}

	return nil, nil
}
