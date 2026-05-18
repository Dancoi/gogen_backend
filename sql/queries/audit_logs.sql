-- name: CreateAuditLog :one
INSERT INTO audit_logs (user_id, action, resource_type, resource_id, changes, status_code, error_message, ip_address, user_agent)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id, user_id, action, resource_type, resource_id, changes, status_code, error_message, ip_address, user_agent, created_at;

-- name: GetUserAuditLogs :many
SELECT id, user_id, action, resource_type, resource_id, changes, status_code, error_message, ip_address, user_agent, created_at
FROM audit_logs
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListAuditLogs :many
SELECT id, user_id, action, resource_type, resource_id, changes, status_code, error_message, ip_address, user_agent, created_at
FROM audit_logs
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetAuditLogsByAction :many
SELECT id, user_id, action, resource_type, resource_id, changes, status_code, error_message, ip_address, user_agent, created_at
FROM audit_logs
WHERE action = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetAuditLogsInDateRange :many
SELECT id, user_id, action, resource_type, resource_id, changes, status_code, error_message, ip_address, user_agent, created_at
FROM audit_logs
WHERE created_at >= $1 AND created_at <= $2
ORDER BY created_at DESC;
