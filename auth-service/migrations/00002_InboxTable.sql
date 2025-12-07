-- +goose Up
CREATE TABLE IF NOT EXISTS inbox (
    id BIGSERIAL PRIMARY KEY,
    message_type TEXT NOT NULL,
    payload JSONB NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
)

CREATE INDEX idx_inbox_status ON inbox(status);
CREATE INDEX idx_created_at ON inbox(created_at);

CREATE TYPE inbox_message_v1 AS (
    id BIGINT,
    message_type TEXT,
    payload JSONB,
    status TEXT,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE
);

-- +goose Down
DROP INDEX IF EXISTS idx_inbox_status;
DROP INDEX IF EXISTS idx_created_at;
DROP TABLE IF EXISTS inbox;
DROP TYPE IF EXISTS inbox_message_v1;
