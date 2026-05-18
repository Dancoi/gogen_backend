-- name: CreateSubscriptionPlan :one
INSERT INTO subscription_plans (name, plan_type, description, max_tokens_per_month, max_api_calls_per_day, price, features, trial_duration_days, is_active)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id, name, plan_type, description, max_tokens_per_month, max_api_calls_per_day, price, features, trial_duration_days, is_active, created_at, updated_at;

-- name: GetSubscriptionPlanByID :one
SELECT id, name, plan_type, description, max_tokens_per_month, max_api_calls_per_day, price, features, trial_duration_days, is_active, created_at, updated_at
FROM subscription_plans
WHERE id = $1;

-- name: GetSubscriptionPlanByType :one
SELECT id, name, plan_type, description, max_tokens_per_month, max_api_calls_per_day, price, features, trial_duration_days, is_active, created_at, updated_at
FROM subscription_plans
WHERE plan_type = $1;

-- name: ListSubscriptionPlans :many
SELECT id, name, plan_type, description, max_tokens_per_month, max_api_calls_per_day, price, features, trial_duration_days, is_active, created_at, updated_at
FROM subscription_plans
WHERE is_active = true
ORDER BY id;

-- name: UpdateSubscriptionPlan :one
UPDATE subscription_plans
SET name = $2, description = $3, max_tokens_per_month = $4, max_api_calls_per_day = $5, price = $6, features = $7, updated_at = NOW()
WHERE id = $1
RETURNING id, name, plan_type, description, max_tokens_per_month, max_api_calls_per_day, price, features, trial_duration_days, is_active, created_at, updated_at;
