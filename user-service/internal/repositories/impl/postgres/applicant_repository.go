package postgresimpl

import (
	"context"
	"fmt"
	"strings"

	"github.com/ZaiiiRan/job_search_service/user-service/internal/domain/user/applicant"
	"github.com/ZaiiiRan/job_search_service/user-service/internal/repositories/interfaces"
	"github.com/ZaiiiRan/job_search_service/user-service/internal/repositories/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ApplicantRepository struct {
	conn *pgxpool.Conn
}

func NewApplicantRepository(conn *pgxpool.Conn) interfaces.ApplicantRepository {
	return &ApplicantRepository{
		conn: conn,
	}
}

func (r *ApplicantRepository) Create(ctx context.Context, applicant *applicant.Applicant) error {
	dal := models.V1ApplicantDalFromDomain(applicant)

	sql := `
		INSERT INTO applicants (
			first_name,
			last_name,
			patronymic,
			birth_date,
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
			(i).first_name,
			(i).last_name,
			(i).patronymic,
			(i).birth_date,
			(i).city,
			(i).email,
			(i).phone_number,
			(i).telegram,
			(i).is_active,
			(i).is_deleted,
			(i).created_at,
			(i).updated_at
		FROM UNNEST($1::v1_applicant[]) i
		RETURNING
			id,
			first_name,
			last_name,
			patronymic,
			birth_date,
			city,
			email,
			phone_number,
			telegram,
			is_active,
			is_deleted,
			created_at,
			updated_at
	`

	rows, err := r.conn.Query(ctx, sql, []models.V1ApplicantDal{dal})
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		var res models.V1ApplicantDal
		if err := rows.Scan(
			&res.Id,
			&res.FirstName,
			&res.LastName,
			&res.Patronymic,
			&res.BirthDate,
			&res.City,
			&res.Email,
			&res.PhoneNumber,
			&res.Telegram,
			&res.IsActive,
			&res.IsDeleted,
			&res.CreatedAt,
			&res.UpdatedAt,
		); err != nil {
			return fmt.Errorf("scan applicant: %w", err)
		}
		applicant = res.ToDomain()
		return nil
	}

	return fmt.Errorf("no applicant returned from insert")
}

func (r *ApplicantRepository) Update(ctx context.Context, applicant *applicant.Applicant) error {
	dal := models.V1ApplicantDalFromDomain(applicant)

	sql := `
		UPDATE applicants AS t
		SET
			first_name   = u.first_name,
			last_name    = u.last_name,
			patronymic   = u.patronymic,
			birth_date   = u.birth_date,
			city         = u.city,
			email        = u.email,
			phone_number = u.phone_number,
			telegram     = u.telegram,
			is_active    = u.is_active,
			is_deleted   = u.is_deleted,
			updated_at   = u.updated_at
		FROM unnest($1::v1_applicant[]) AS u
		WHERE t.id = u.id
		RETURNING
			t.id,
			t.first_name,
			t.last_name,
			t.patronymic,
			t.birth_date,
			t.city,
			t.email,
			t.phone_number,
			t.telegram,
			t.is_active,
			t.is_deleted,
			t.created_at,
			t.updated_at;
	`

	rows, err := r.conn.Query(ctx, sql, []models.V1ApplicantDal{dal})
	if err != nil {
		return fmt.Errorf("update applicant: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		var res models.V1ApplicantDal
		if err := rows.Scan(
			&res.Id,
			&res.FirstName,
			&res.LastName,
			&res.Patronymic,
			&res.BirthDate,
			&res.City,
			&res.Email,
			&res.PhoneNumber,
			&res.Telegram,
			&res.IsActive,
			&res.IsDeleted,
			&res.CreatedAt,
			&res.UpdatedAt,
		); err != nil {
			return fmt.Errorf("scan applicant: %w", err)
		}
		applicant = res.ToDomain()
		return nil
	}

	return fmt.Errorf("no applicant updated")
}

func (r *ApplicantRepository) Query(ctx context.Context, query *models.QueryApplicantsDal) ([]*applicant.Applicant, error) {
	if query == nil {
		query = &models.QueryApplicantsDal{}
	}

	var (
		sb     strings.Builder
		args   []any
		argPos = 1
	)

	sb.WriteString(`
		SELECT
			id, first_name, last_name, patronymic, birth_date, city,
			email, phone_number, telegram, is_active, is_deleted,
			created_at, updated_at
		FROM applicants
		WHERE 1=1
	`)

	appendAnyEqual(&sb, "id", query.Ids, &args, &argPos)
	appendAnyEqual(&sb, "email", query.Emails, &args, &argPos)
	appendILike(&sb, "email", query.EmailSubstrs, &args, &argPos)
	appendBool(&sb, "is_active", query.IsActive, &args, &argPos)
	appendBool(&sb, "is_deleted", query.IsDeleted, &args, &argPos)
	appendOrder(&sb, "id", true)
	appendLimitOffset(&sb, query.Limit, query.Offset, &args, &argPos)

	rows, err := r.conn.Query(ctx, sb.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("query applicants: %w", err)
	}
	defer rows.Close()

	var result []*applicant.Applicant
	for rows.Next() {
		var dal models.V1ApplicantDal
		if err := rows.Scan(
			&dal.Id, &dal.FirstName, &dal.LastName, &dal.Patronymic, &dal.BirthDate, &dal.City,
			&dal.Email, &dal.PhoneNumber, &dal.Telegram, &dal.IsActive, &dal.IsDeleted,
			&dal.CreatedAt, &dal.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan applicant: %w", err)
		}
		result = append(result, dal.ToDomain())
	}

	return result, nil
}
