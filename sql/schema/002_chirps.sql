-- +goose Up
CREATE TABLE chirps (
    id uuid PRIMARY KEY,
    body text not null,
    created_at timestamp not null,
    user_id uuid not null,
    updated_at timestamp not null,
    CONSTRAINT user_constraint FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE chirps;