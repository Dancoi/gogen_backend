-- name: CreateUserSession :one
INSERT INTO user_sessions (user_id, session_token, ip_address, user_agent, expires_at, last_activity_at)
VALUES ($1, $2, $3, $4, $5, NOW())
RETURNING id, user_id, session_token, ip_address, user_agent, expires_at, last_activity_at, created_at, updated_at;

-- name: GetUserSessionByToken :one
SELECT id, user_id, session_token, ip_address, user_agent, expires_at, last_activity_at, created_at, updated_at
FROM user_sessions
WHERE session_token = $1 AND expires_at > NOW();

-- name: GetUserSessionByID :one
SELECT id, user_id, session_token, ip_address, user_agent, expires_at, last_activity_at, created_at, updated_at
FROM user_sessions
WHERE id = $1;

-- name: ListUserSessions :many
SELECT id, user_id, session_token, ip_address, user_agent, expires_at, last_activity_at, created_at, updated_at
FROM user_sessions
WHERE user_id = $1 AND expires_at > NOW()
ORDER BY created_at DESC;

-- name: UpdateUserSessionActivity :one
UPDATE user_sessions
SET last_activity_at = NOW(), updated_at = NOW()
WHERE id = $1
RETURNING id, user_id, session_token, ip_address, user_agent, expires_at, last_activity_at, created_at, updated_at;

-- name: RevokeUserSession :exec
DELETE FROM user_sessions
WHERE id = $1;

-- name: RevokeAllUserSessions :exec
DELETE FROM user_sessions
WHERE user_id = $1;

-- name: DeleteExpiredSessions :exec
DELETE FROM user_sessions
WHERE expires_at < NOW();

-- name: CountActiveUserSessions :one
SELECT COUNT(*) FROM user_sessions
WHERE user_id = $1 AND expires_at > NOW();
