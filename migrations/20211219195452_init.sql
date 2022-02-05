-- +goose Up
-- +goose StatementBegin
CREATE TABLE estates
(
    id         SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT (now()),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,

    ext_id VARCHAR NOT NULL,
    title       VARCHAR NOT NULL,
    url         VARCHAR NOT NULL,
    price       VARCHAR NOT NULL
);
COMMENT ON TABLE estates IS 'Estates taken from estate providers';
COMMENT ON COLUMN estates.ext_id IS 'Estate provider estate_id';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE estates;
-- +goose StatementEnd
