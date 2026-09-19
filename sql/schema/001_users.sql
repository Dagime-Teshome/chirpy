-- +goose Up

CREATE TABLE users (
    id uuid PRIMARY KEY,
    email text UNIQUE NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL
);

-- +goose Down

DROP TABLE users;