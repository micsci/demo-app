-- +goose Up

ALTER TABLE application ADD COLUMN price DOUBLE PRECISION NOT NULL DEFAULT 0;

CREATE TABLE review (
    id             UUID PRIMARY KEY,
    application_id UUID NOT NULL REFERENCES application(id),
    author_email   TEXT NOT NULL,
    rating         INTEGER NOT NULL,
    comment        TEXT NOT NULL DEFAULT '',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TRIGGER set_updated_at BEFORE UPDATE ON review
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();
