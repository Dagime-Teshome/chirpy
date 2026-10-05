-- name: ListUsers :many
SELECT * FROM users 
ORDER BY created_at;

-- name: GetUserByID :one
SELECT * FROM users 
WHERE id = $1;


-- name: GerUserByEmail :one
SELECT * FROM users
WHERE email = $1;


-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email , hashed_password)
VALUES (
    gen_random_uuid(),NOW(),NOW(),$1,$2
)
RETURNING *;

-- name: UpdateUser :one
UPDATE users 
SET email = $1,
    hashed_password = $2,
    updated_at = NOw()

WHERE id = $3
RETURNING *;

-- name: UpdateUserMembership :one
UPDATE users 
SET is_chirpy_red = $1,
    updated_at = NOw()
WHERE id = $2
RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users;