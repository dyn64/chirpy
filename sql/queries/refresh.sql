-- name: AddRefreshToken :exec
INSERT INTO refresh_token(token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES(
$1,
NOW(),
NOW(),
$2,
$3,
NULL);

-- name: bad_GetUserFromRefreshToken_bad :one
SELECT * FROM refresh_token
WHERE token = $1;

-- name: GetUserFromRefreshToken :one
SELECT users.* FROM users
JOIN refresh_token ON users.id = refresh_token.user_id
WHERE refresh_token.token = $1
AND revoked_at is NULL
AND expires_at > NOW();

-- name: RevokeRefreshToken :exec
UPDATE refresh_token
SET revoked_at = NOW()::TIMESTAMP,
updated_at = NOW()::TIMESTAMP
where token = $1;
