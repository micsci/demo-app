-- +goose Up

CREATE TABLE review (
    id             UUID PRIMARY KEY,
    application_id UUID NOT NULL REFERENCES application(id) ON DELETE CASCADE,
    merchant_id    TEXT NOT NULL,
    rating         INTEGER NOT NULL,
    body           TEXT NOT NULL DEFAULT '',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (application_id, merchant_id)
);

-- +goose Down
DROP TABLE IF EXISTS review;
