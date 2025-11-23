-- +goose Up
CREATE TABLE IF NOT EXISTS applicant_passwords (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE UNIQUE INDEX idx_applicant_passwords_user_id ON applicant_passwords (user_id);

CREATE TABLE applicant_activation_codes (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE,
    code TEXT NOT NULL,
    generations_left INT NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE UNIQUE INDEX idx_applicant_activation_codes_user_id ON applicant_activation_codes (user_id);

CREATE TABLE applicant_reset_password_codes (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE,
    code TEXT NOT NULL,
    generations_left INT NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE UNIQUE INDEX idx_applicant_reset_password_codes_user_id ON applicant_reset_password_codes(user_id);

CREATE TABLE applicant_refresh_tokens (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    token TEXT NOT NULL,
    version INTEGER NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_applicant_refresh_tokens_user_id ON applicant_refresh_tokens(user_id);
CREATE INDEX idx_applicant_refresh_tokens_token ON applicant_refresh_tokens(token);

CREATE TABLE applicant_version (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    version INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE UNIQUE INDEX idx_applicant_version_user_id ON applicant_version(user_id);

CREATE TABLE IF NOT EXISTS employer_passwords (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE UNIQUE INDEX idx_employer_passwords_user_id ON employer_passwords (user_id);

CREATE TABLE employer_activation_codes (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE,
    code TEXT NOT NULL,
    generations_left INT NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE UNIQUE INDEX idx_employer_activation_codes_user_id ON employer_activation_codes (user_id);

CREATE TABLE employer_reset_password_codes (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE,
    code TEXT NOT NULL,
    generations_left INT NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE UNIQUE INDEX idx_employer_reset_password_codes_user_id ON employer_reset_password_codes(user_id);

CREATE TABLE employer_refresh_tokens (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    token TEXT NOT NULL,
    version INTEGER NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_employer_refresh_tokens_user_id ON employer_refresh_tokens(user_id);
CREATE INDEX idx_employer_refresh_tokens_token ON employer_refresh_tokens(token);

CREATE TABLE employer_version (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    version INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE UNIQUE INDEX idx_employer_version_user_id ON employer_version(user_id);

CREATE TYPE v1_user_password AS (
    id BIGINT,
    user_id BIGINT,
    password_hash TEXT,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE
);

CREATE TYPE v1_code AS (
    id BIGINT,
    user_id BIGINT,
    code TEXT,
    generations_left INT,
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE
);

CREATE TYPE v1_refresh_token AS (
    id BIGINT,
    user_id BIGINT,
    token TEXT,
    version INTEGER,
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE
);

CREATE TYPE v1_user_version AS (
    id BIGINT,
    user_id BIGINT,
    version INTEGER,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE
);

-- +goose Down
DROP INDEX IF EXISTS idx_applicant_passwords_user_id;
DROP TABLE IF EXISTS applicant_passwords;

DROP INDEX IF EXISTS idx_applicant_activation_codes_user_id;
DROP TABLE IF EXISTS applicant_activation_codes;

DROP INDEX IF EXISTS idx_applicant_reset_password_codes_user_id;
DROP TABLE IF EXISTS applicant_reset_password_codes;

DROP INDEX IF EXISTS idx_applicant_refresh_tokens_user_id;
DROP INDEX IF EXISTS idx_applicant_refresh_tokens_token;
DROP TABLE IF EXISTS applicant_refresh_tokens;

DROP INDEX IF EXISTS idx_applicant_version_user_id;
DROP TABLE IF EXISTS applicant_version;

DROP INDEX IF EXISTS idx_employer_passwords_user_id;
DROP TABLE IF EXISTS employer_passwords;

DROP INDEX IF EXISTS idx_employer_activation_codes_user_id;
DROP TABLE IF EXISTS employer_activation_codes;

DROP INDEX IF EXISTS idx_employer_reset_password_codes_user_id;
DROP TABLE IF EXISTS employer_reset_password_codes;

DROP INDEX IF EXISTS idx_employer_refresh_tokens_user_id;
DROP INDEX IF EXISTS idx_employer_refresh_tokens_token;
DROP TABLE IF EXISTS employer_refresh_tokens;

DROP INDEX IF EXISTS idx_employer_version_user_id;
DROP TABLE IF EXISTS employer_version;

DROP TYPE IF EXISTS v1_user_password;
DROP TYPE IF EXISTS v1_code;
DROP TYPE IF EXISTS v1_refresh_token;
DROP TYPE IF EXISTS v1_user_version;
