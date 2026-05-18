-- name: CreateApiToken :one
INSERT INTO api_tokens (user_id, subscription_id, name, token_hash, is_active)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, user_id, subscription_id, name, token_hash, is_active, created_at, last_used_at, expires_at, revoked_at;

-- name: GetApiTokenByID :one
SELECT id, user_id, subscription_id, name, token_hash, is_active, created_at, last_used_at, expires_at, revoked_at
FROM api_tokens
WHERE id = $1;

-- name: GetApiTokenByHash :one
SELECT id, user_id, subscription_id, name, token_hash, is_active, created_at, last_used_at, expires_at, revoked_at
FROM api_tokens
WHERE token_hash = $1;

-- name: ListUserApiTokens :many
SELECT id, user_id, subscription_id, name, token_hash, is_active, created_at, last_used_at, expires_at, revoked_at
FROM api_tokens
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateApiTokenLastUsed :one
UPDATE api_tokens
SET last_used_at = NOW()
WHERE id = $1
RETURNING id, user_id, subscription_id, name, token_hash, is_active, created_at, last_used_at, expires_at, revoked_at;

-- name: RevokeApiToken :one
UPDATE api_tokens
SET is_active = false, revoked_at = NOW()
WHERE id = $1
RETURNING id, user_id, subscription_id, name, token_hash, is_active, created_at, last_used_at, expires_at, revoked_at;

-- name: DeleteApiToken :exec
DELETE FROM api_tokens
WHERE id = $1;

-- name: CountUserActiveTokens :one
SELECT COUNT(*) FROM api_tokens
WHERE user_id = $1 AND is_active = true AND (expires_at IS NULL OR expires_at > NOW()) AND revoked_at IS NULL;

-- name: ListExpiredApiTokens :many
SELECT id, user_id, subscription_id, name, token_hash, is_active, created_at, last_used_at, expires_at, revoked_at
FROM api_tokens
WHERE expires_at IS NOT NULL AND expires_at < NOW() AND is_active = true;
