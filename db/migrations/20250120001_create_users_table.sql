-- +goose Up
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    active BOOLEAN NOT NULL
);

-- +goose Down
DROP TABLE users;