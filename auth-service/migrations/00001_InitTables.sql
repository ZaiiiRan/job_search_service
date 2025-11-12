-- +goose Up
CREATE TABLE IF NOT EXISTS user_passwords (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE UNIQUE INDEX idx_user_passwords_user_id ON user_passwords (user_id);

CREATE TYPE v1_user_password AS (
    id BIGINT,
    user_id BIGINT,
    password_hash TEXT,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE activation_codes (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE,
    code TEXT NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE UNIQUE INDEX idx_activation_codes_user_id ON activation_codes (user_id);

CREATE TYPE v1_activation_code AS (
    id BIGINT,
    user_id BIGINT,
    code TEXT,
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE reset_password_codes (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE,
    code TEXT NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE UNIQUE INDEX idx_reset_password_codes_user_id ON reset_password_codes(user_id);

CREATE TYPE v1_reset_password_code AS (
    id BIGINT,
    user_id BIGINT,
    code TEXT,
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE refresh_tokens (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    token TEXT NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_token ON refresh_tokens(token);

CREATE TYPE v1_refresh_token AS (
    id BIGINT,
    user_id BIGINT,
    token TEXT,
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE
);

-- +goose Down
DROP INDEX IF EXISTS idx_refresh_tokens_token;
DROP INDEX IF EXISTS idx_refresh_tokens_user_id;
DROP TABLE IF EXISTS refresh_tokens;

DROP INDEX IF EXISTS idx_reset_password_codes_user_id;
DROP TABLE IF EXISTS reset_password_codes;

DROP INDEX IF EXISTS idx_activation_codes_user_id;
DROP TABLE IF EXISTS activation_codes;

DROP INDEX IF EXISTS idx_user_passwords_user_id;
DROP TABLE IF EXISTS user_passwords;

DROP TYPE IF EXISTS v1_refresh_token;
DROP TYPE IF EXISTS v1_reset_password_code;
DROP TYPE IF EXISTS v1_activation_code;
DROP TYPE IF EXISTS v1_user_password;
