-- +goose Up
INSERT INTO roles (name) VALUES ('admin'), ('user');
