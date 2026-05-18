-- name: CreateTokenUsageLog :one
INSERT INTO token_usage_logs (api_token_id, user_id, endpoint, method, status_code, response_time, data_size, error_message, ip_address, user_agent)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING id, api_token_id, user_id, endpoint, method, status_code, response_time, data_size, error_message, ip_address, user_agent, created_at;

-- name: GetTokenUsageLogs :many
SELECT id, api_token_id, user_id, endpoint, method, status_code, response_time, data_size, error_message, ip_address, user_agent, created_at
FROM token_usage_logs
WHERE api_token_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetUserTokenUsageLogs :many
SELECT id, api_token_id, user_id, endpoint, method, status_code, response_time, data_size, error_message, ip_address, user_agent, created_at
FROM token_usage_logs
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetTokenUsageLogsInDateRange :many
SELECT id, api_token_id, user_id, endpoint, method, status_code, response_time, data_size, error_message, ip_address, user_agent, created_at
FROM token_usage_logs
WHERE api_token_id = $1 AND created_at >= $2 AND created_at <= $3
ORDER BY created_at DESC;

-- name: CountSuccessfulRequests :one
SELECT COUNT(*) FROM token_usage_logs
WHERE api_token_id = $1 AND status_code >= 200 AND status_code < 300;

-- name: CountFailedRequests :one
SELECT COUNT(*) FROM token_usage_logs
WHERE api_token_id = $1 AND (status_code < 200 OR status_code >= 300);

-- name: GetAverageResponseTime :one
SELECT AVG(response_time) as avg_response_time FROM token_usage_logs
WHERE api_token_id = $1;
