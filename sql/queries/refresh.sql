-- name: AddRefreshToken :exec
INSERT INTO refresh_token(token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES(
$1,
NOW(),
NOW(),
$2,
$3,
NULL);

-- name: GetUserFromRefreshToken :one
SELECT * FROM refresh_token
WHERE token = $1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_token
SET revoked_at = NOW()::TIMESTAMP,
updated_at = NOW()::TIMESTAMP
where token = $1;
