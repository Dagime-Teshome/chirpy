-- +goose Up
CREATE TABLE refresh_tokens(
    token text not null PRIMARY key,
    created_at timestamp not null ,
    updated_at timestamp not null,
    expires_at timestamp not null,
    revoked_at timestamp default null,
    user_id uuid not null REFERENCES users(id) ON DELETE CASCADE
);




-- +goose Down
DROP TABLE refresh_tokens;