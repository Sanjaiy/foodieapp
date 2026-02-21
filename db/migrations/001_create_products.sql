-- +goose Up
CREATE TABLE IF NOT EXISTS products (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    price       NUMERIC(10,2) NOT NULL,
    category    TEXT NOT NULL,
    img_thumb   TEXT NOT NULL DEFAULT '',
    img_mobile  TEXT NOT NULL DEFAULT '',
    img_tablet  TEXT NOT NULL DEFAULT '',
    img_desktop TEXT NOT NULL DEFAULT ''
);

-- +goose Down
DROP TABLE IF EXISTS products;
