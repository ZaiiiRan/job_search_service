package postgresimpl

import (
	"context"
	"fmt"
	"strings"

	"github.com/ZaiiiRan/job_search_service/user-service/internal/domain/user/employer"
	"github.com/ZaiiiRan/job_search_service/user-service/internal/repositories/interfaces"
	"github.com/ZaiiiRan/job_search_service/user-service/internal/repositories/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EmployerRepository struct {
	conn *pgxpool.Conn
}

func NewEmployerRepository(conn *pgxpool.Conn) interfaces.EmployerRepository {
	return &EmployerRepository{
		conn: conn,
	}
}

func (r *EmployerRepository) Create(ctx context.Context, e *employer.Employer) error {
	dal := models.V1EmployerDalFromDomain(e)

	sql := `
		INSERT INTO employers (
			company_name,
			city,
			email,
			phone_number,
			telegram,
			is_active,
			is_deleted,
			created_at,
			updated_at
		)
		SELECT
			(i).company_name,
			(i).city,
			(i).email,
			(i).phone_number,
			(i).telegram,
			(i).is_active,
			(i).is_deleted,
			(i).created_at,
			(i).updated_at
		FROM UNNEST($1::v1_employer[]) AS i
		RETURNING
			id,
			company_name,
			city,
			email,
			phone_number,
			telegram,
			is_active,
			is_deleted,
			created_at,
			updated_at;
	`

	rows, err := r.conn.Query(ctx, sql, []models.V1EmployerDal{dal})
	if err != nil {
		return fmt.Errorf("insert employer: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		var res models.V1EmployerDal
		if err := rows.Scan(
			&res.Id,
			&res.CompanyName,
			&res.City,
			&res.Email,
			&res.PhoneNumber,
			&res.Telegram,
			&res.IsActive,
			&res.IsDeleted,
			&res.CreatedAt,
			&res.UpdatedAt,
		); err != nil {
			return fmt.Errorf("scan employer: %w", err)
		}
		e = res.ToDomain()
		return nil
	}

	return fmt.Errorf("no employer returned from insert")
}

func (r *EmployerRepository) Update(ctx context.Context, e *employer.Employer) error {
	dal := models.V1EmployerDalFromDomain(e)

	sql := `
		UPDATE employers AS t
		SET
			company_name = u.company_name,
			city         = u.city,
			email        = u.email,
			phone_number = u.phone_number,
			telegram     = u.telegram,
			is_active    = u.is_active,
			is_deleted   = u.is_deleted,
			updated_at   = u.updated_at
		FROM UNNEST($1::v1_employer[]) AS u
		WHERE t.id = u.id
		RETURNING
			t.id,
			t.company_name,
			t.city,
			t.email,
			t.phone_number,
			t.telegram,
			t.is_active,
			t.is_deleted,
			t.created_at,
			t.updated_at;
	`

	rows, err := r.conn.Query(ctx, sql, []models.V1EmployerDal{dal})
	if err != nil {
		return fmt.Errorf("update employer: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		var res models.V1EmployerDal
		if err := rows.Scan(
			&res.Id,
			&res.CompanyName,
			&res.City,
			&res.Email,
			&res.PhoneNumber,
			&res.Telegram,
			&res.IsActive,
			&res.IsDeleted,
			&res.CreatedAt,
			&res.UpdatedAt,
		); err != nil {
			return fmt.Errorf("scan employer: %w", err)
		}
		e = res.ToDomain()
		return nil
	}

	return fmt.Errorf("no employer updated")
}

func (r *EmployerRepository) Query(ctx context.Context, q *models.QueryEmployersDal) ([]*employer.Employer, error) {
	if q == nil {
		q = &models.QueryEmployersDal{}
	}

	var (
		sb     strings.Builder
		args   []any
		argPos = 1
	)

	sb.WriteString(`
		SELECT
			id,
			company_name,
			city,
			email,
			phone_number,
			telegram,
			is_active,
			is_deleted,
			created_at,
			updated_at
		FROM employers
		WHERE 1=1
	`)

	appendAnyEqual(&sb, "id", q.Ids, &args, &argPos)
	appendAnyEqual(&sb, "email", q.Emails, &args, &argPos)
	appendILike(&sb, "email", q.EmailSubstrs, &args, &argPos)
	appendAnyEqual(&sb, "company_name", q.CompanyNames, &args, &argPos)
	appendILike(&sb, "company_name", q.CompanyNameSubstrs, &args, &argPos)
	appendBool(&sb, "is_active", q.IsActive, &args, &argPos)
	appendBool(&sb, "is_deleted", q.IsDeleted, &args, &argPos)
	appendOrder(&sb, "id", true)
	appendLimitOffset(&sb, q.Limit, q.Offset, &args, &argPos)

	rows, err := r.conn.Query(ctx, sb.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("query employers: %w", err)
	}
	defer rows.Close()

	var result []*employer.Employer
	for rows.Next() {
		var dal models.V1EmployerDal
		if err := rows.Scan(
			&dal.Id,
			&dal.CompanyName,
			&dal.City,
			&dal.Email,
			&dal.PhoneNumber,
			&dal.Telegram,
			&dal.IsActive,
			&dal.IsDeleted,
			&dal.CreatedAt,
			&dal.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan employer: %w", err)
		}
		result = append(result, dal.ToDomain())
	}

	return result, nil
}
