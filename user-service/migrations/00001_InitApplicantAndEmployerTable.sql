-- +goose Up
CREATE TABLE IF NOT EXISTS applicants (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    patronymic TEXT,
    birth_date DATE NOT NULL,
    city TEXT NOT NULL,
    email TEXT NOT NULL,
    phone_number TEXT,
    telegram TEXT,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE UNIQUE INDEX idx_applicants_email_active_unique
    ON applicants (email)
    WHERE is_deleted = FALSE;

CREATE TYPE v1_applicant as (
    id BIGINT,
    first_name TEXT,
    last_name TEXT,
    patronymic TEXT,
    birth_date DATE,
    city TEXT,
    email TEXT,
    phone_number TEXT,
    telegram TEXT,
    is_active BOOLEAN,
    is_deleted BOOLEAN,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS employers (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    company_name TEXT NOT NULL,
    city TEXT NOT NULL,
    email TEXT NOT NULL,
    phone_number TEXT,
    telegram TEXT,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE UNIQUE INDEX idx_employers_email_active_unique
    ON employers (email)
    WHERE is_deleted = FALSE;

CREATE TYPE v1_employer AS (
    id BIGINT,
    company_name TEXT,
    city TEXT,
    email TEXT,
    phone_number TEXT,
    telegram TEXT,
    is_active BOOLEAN,
    is_deleted BOOLEAN,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE
);

-- +goose Down
DROP INDEX IF EXISTS idx_applicants_email_active_unique;
DROP INDEX IF EXISTS idx_employers_email_active_unique;
DROP TABLE IF EXISTS applicants;
DROP TYPE IF EXISTS v1_applicant;
DROP TABLE IF EXISTS employer;
DROP TYPE IF EXISTS v1_employer;
