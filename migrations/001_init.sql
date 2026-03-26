-- +goose Up

CREATE TABLE category (
    id         UUID PRIMARY KEY,
    name       TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE application (
    id           UUID PRIMARY KEY,
    name         TEXT NOT NULL UNIQUE,
    display_name TEXT NOT NULL DEFAULT '',
    description  TEXT NOT NULL DEFAULT '',
    active       BOOLEAN NOT NULL DEFAULT true,
    visible      BOOLEAN NOT NULL DEFAULT true,
    category_id  UUID NOT NULL REFERENCES category(id),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE connection (
    id             UUID PRIMARY KEY,
    application_id UUID NOT NULL REFERENCES application(id) ON DELETE CASCADE,
    merchant_id    TEXT NOT NULL,
    status         TEXT NOT NULL DEFAULT 'active',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (application_id, merchant_id)
);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION trigger_set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER set_updated_at BEFORE UPDATE ON category
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

CREATE TRIGGER set_updated_at BEFORE UPDATE ON application
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

CREATE TRIGGER set_updated_at BEFORE UPDATE ON connection
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS connection;
DROP TABLE IF EXISTS application;
DROP TABLE IF EXISTS category;
DROP FUNCTION IF EXISTS trigger_set_updated_at;
