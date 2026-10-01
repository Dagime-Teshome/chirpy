-- name: Create_Refresh_Token :one
INSERT INTO refresh_tokens (token , user_id, created_at ,updated_at ,expires_at , revoked_at)
VALUES(
    $1,$2,NOW(),NOW(),NOW() + INTERVAL '60 days',NULL
)
RETURNING *;

-- name: GetUserFromRefreshToken :one
SELECT * from refresh_tokens 
where token = $1;


-- name: Revoke_Refresh_Token :exec
UPDATE refresh_tokens
SET revoked_at = now(),
    updated_at = now()
WHERE token = $1 AND revoked_at IS NULL;;