package postgresimpl

import (
	"context"
	"fmt"
	"strings"

	"github.com/ZaiiiRan/job_search_service/auth-service/internal/domain/code"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/impl"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/interfaces"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/models"
	postgresunitofwork "github.com/ZaiiiRan/job_search_service/auth-service/internal/repositories/unitofwork/postgres"
)

const (
	ApplicantActivationCodes    impl.RepositoryType = "applicant_activation_codes"
	EmployerActivationCodes     impl.RepositoryType = "employer_activation_codes"
	ApplicantResetPasswordCodes impl.RepositoryType = "applicant_reset_password_codes"
	EmployerResetPasswordCodes  impl.RepositoryType = "employer_reset_password_codes"
)

type CodeRepository struct {
	uow       *postgresunitofwork.UnitOfWork
	tableName string
}

func NewCodeRepository(uow *postgresunitofwork.UnitOfWork, repoType impl.RepositoryType) interfaces.CodeRepository {
	return &CodeRepository{uow: uow, tableName: string(repoType)}
}

func (r *CodeRepository) CreateCode(ctx context.Context, code *code.Code) error {
	dal := models.V1CodeDalFromDomain(code)

	conn, err := r.uow.GetConn(ctx)
	if err != nil {
		return err
	}

	var sb strings.Builder

	sb.WriteString(`
		INSERT INTO ` + r.tableName + ` (
			user_id,
			code,
			generations_left,
			expires_at,
			created_at,
			updated_at
		)
		SELECT
			(i).user_id,
			(i).code,
			(i).generations_left,
			(i).expires_at,
			(i).created_at,
			(i).updated_at
		FROM UNNEST($1::v1_code[]) i
		RETURNING
			id,
			user_id,
			code,
			generations_left,
			expires_at,
			created_at,
			updated_at
	`)

	rows, err := conn.Query(ctx, sb.String(), []models.V1CodeDal{dal})
	if err != nil {
		return fmt.Errorf("insert code: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		var res models.V1CodeDal
		if err := rows.Scan(
			&res.Id,
			&res.UserId,
			&res.Code,
			&res.GenerationsLeft,
			&res.ExpiresAt,
			&res.CreatedAt,
			&res.UpdatedAt,
		); err != nil {
			return fmt.Errorf("scan code: %w", err)
		}
		*code = *res.ToDomain()
		return nil
	}

	return fmt.Errorf("no code returned from insert")
}

func (r *CodeRepository) UpdateCode(ctx context.Context, code *code.Code) error {
	dal := models.V1CodeDalFromDomain(code)

	conn, err := r.uow.GetConn(ctx)
	if err != nil {
		return err
	}

	var sb strings.Builder

	sb.WriteString(`
		UPDATE ` + r.tableName + ` AS t
		SET
			user_id = (i).user_id,
			code = (i).code,
			generations_left = (i).generations_left,
			expires_at = (i).expires_at,
			created_at = (i).created_at,
			updated_at = (i).updated_at
		FROM UNNEST($1::v1_code[]) i
		WHERE t.id = (i).id
		RETURNING
			t.id,
			t.user_id,
			t.code,
			t.generations_left,
			t.expires_at,
			t.created_at,
			t.updated_at
	`)

	rows, err := conn.Query(ctx, sb.String(), []models.V1CodeDal{dal})
	if err != nil {
		return fmt.Errorf("update code: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		var res models.V1CodeDal
		if err := rows.Scan(
			&res.Id,
			&res.UserId,
			&res.Code,
			&res.GenerationsLeft,
			&res.ExpiresAt,
			&res.CreatedAt,
			&res.UpdatedAt,
		); err != nil {
			return fmt.Errorf("scan code: %w", err)
		}
		*code = *res.ToDomain()
		return nil
	}

	return fmt.Errorf("no code updated")
}

func (r *CodeRepository) DeleteCode(ctx context.Context, code *code.Code) error {
	dal := models.V1CodeDalFromDomain(code)

	conn, err := r.uow.GetConn(ctx)
	if err != nil {
		return err
	}

	var sb strings.Builder

	sb.WriteString(`
		DELETE FROM ` + r.tableName + ` AS t
		USING (
			SELECT (i).id AS id
			FROM UNNEST($1::v1_code[]) AS i
		) AS d
		WHERE t.id = d.id
		RETURNING
			t.id,
			t.user_id,
			t.code,
			t.generations_left,
			t.expires_at,
			t.created_at,
			t.updated_at
	`)

	rows, err := conn.Query(ctx, sb.String(), []models.V1CodeDal{dal})
	if err != nil {
		return fmt.Errorf("delete code: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		var res models.V1CodeDal
		if err := rows.Scan(
			&res.Id,
			&res.UserId,
			&res.Code,
			&res.GenerationsLeft,
			&res.ExpiresAt,
			&res.CreatedAt,
			&res.UpdatedAt,
		); err != nil {
			return fmt.Errorf("scan code: %w", err)
		}
		*code = *res.ToDomain()
		return nil
	}

	return fmt.Errorf("no code deleted")
}

func (r *CodeRepository) QueryCode(ctx context.Context, query *models.QueryCodeDal) (*code.Code, error) {
	if query == nil {
		query = &models.QueryCodeDal{}
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
			id, user_id, code,
			generations_left, expires_at,
			created_at, updated_at
		FROM ` + r.tableName + `
		WHERE 1=1
	`)

	appendEqual(&sb, "id", query.Id, &args, &argPos)
	appendEqual(&sb, "user_id", query.UserId, &args, &argPos)
	appendLimitOffset(&sb, 1, 0, &args, &argPos)

	rows, err := conn.Query(ctx, sb.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("query code: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		var res models.V1CodeDal
		if err := rows.Scan(
			&res.Id,
			&res.UserId,
			&res.Code,
			&res.GenerationsLeft,
			&res.ExpiresAt,
			&res.CreatedAt,
			&res.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan code: %w", err)
		}
		code := res.ToDomain()
		return code, nil
	}

	return nil, nil
}
