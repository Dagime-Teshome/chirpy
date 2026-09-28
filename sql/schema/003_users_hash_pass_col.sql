-- +goose Up
ALTER TABLE users
ADD COLUMN hashed_password TEXT  NOT NULL DEFAULT 'not_set';;

-- +goose Down
ALTER TABLE users
DROP COLUMN hashed_password;