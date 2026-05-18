-- name: CreateSubscription :one
INSERT INTO subscriptions (user_id, subscription_plan_id, plan_type, is_active, started_at, expires_at, max_tokens_per_month, current_usage, usage_reset_date)
VALUES ($1, $2, $3, $4, NOW(), $5, $6, 0, NOW())
RETURNING id, user_id, subscription_plan_id, plan_type, is_active, started_at, expires_at, renewal_date, max_tokens_per_month, current_usage, usage_reset_date, created_at, updated_at;

-- name: GetSubscriptionByID :one
SELECT id, user_id, subscription_plan_id, plan_type, is_active, started_at, expires_at, renewal_date, max_tokens_per_month, current_usage, usage_reset_date, created_at, updated_at
FROM subscriptions
WHERE id = $1;

-- name: GetUserActiveSubscription :one
SELECT id, user_id, subscription_plan_id, plan_type, is_active, started_at, expires_at, renewal_date, max_tokens_per_month, current_usage, usage_reset_date, created_at, updated_at
FROM subscriptions
WHERE user_id = $1 AND is_active = true AND expires_at > NOW()
LIMIT 1;

-- name: GetUserSubscription :one
SELECT id, user_id, subscription_plan_id, plan_type, is_active, started_at, expires_at, renewal_date, max_tokens_per_month, current_usage, usage_reset_date, created_at, updated_at
FROM subscriptions
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT 1;

-- name: UpdateSubscriptionCurrentUsage :one
UPDATE subscriptions
SET current_usage = $2, updated_at = NOW()
WHERE id = $1
RETURNING id, user_id, subscription_plan_id, plan_type, is_active, started_at, expires_at, renewal_date, max_tokens_per_month, current_usage, usage_reset_date, created_at, updated_at;

-- name: IncrementSubscriptionUsage :one
UPDATE subscriptions
SET current_usage = current_usage + 1, updated_at = NOW()
WHERE id = $1
RETURNING id, user_id, subscription_plan_id, plan_type, is_active, started_at, expires_at, renewal_date, max_tokens_per_month, current_usage, usage_reset_date, created_at, updated_at;

-- name: UpdateSubscriptionIsActive :one
UPDATE subscriptions
SET is_active = $2, updated_at = NOW()
WHERE id = $1
RETURNING id, user_id, subscription_plan_id, plan_type, is_active, started_at, expires_at, renewal_date, max_tokens_per_month, current_usage, usage_reset_date, created_at, updated_at;

-- name: ResetSubscriptionUsage :one
UPDATE subscriptions
SET current_usage = 0, usage_reset_date = NOW(), updated_at = NOW()
WHERE id = $1
RETURNING id, user_id, subscription_plan_id, plan_type, is_active, started_at, expires_at, renewal_date, max_tokens_per_month, current_usage, usage_reset_date, created_at, updated_at;

-- name: ListUserSubscriptions :many
SELECT id, user_id, subscription_plan_id, plan_type, is_active, started_at, expires_at, renewal_date, max_tokens_per_month, current_usage, usage_reset_date, created_at, updated_at
FROM subscriptions
WHERE user_id = $1
ORDER BY created_at DESC;
