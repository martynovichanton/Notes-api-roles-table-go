-- +goose Up
CREATE TABLE roles (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);
-- +goose Down
DROP TABLE roles;