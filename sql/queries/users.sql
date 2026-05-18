-- name: CreateUser :one
INSERT INTO users (email, username, password_hash, is_active, email_verified)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, email, username, password_hash, is_active, email_verified, created_at, updated_at;

-- name: GetUserByID :one
SELECT id, email, username, password_hash, is_active, email_verified, created_at, updated_at
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, email, username, password_hash, is_active, email_verified, created_at, updated_at
FROM users
WHERE email = $1;

-- name: GetUserByUsername :one
SELECT id, email, username, password_hash, is_active, email_verified, created_at, updated_at
FROM users
WHERE username = $1;

-- name: ListUsers :many
SELECT id, email, username, is_active, email_verified, created_at, updated_at
FROM users
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET username = $2, email = $3, updated_at = NOW()
WHERE id = $1
RETURNING id, email, username, password_hash, is_active, email_verified, created_at, updated_at;

-- name: UpdateUserPasswordHash :one
UPDATE users
SET password_hash = $2, updated_at = NOW()
WHERE id = $1
RETURNING id, email, username, password_hash, is_active, email_verified, created_at, updated_at;

-- name: UpdateUserIsActive :one
UPDATE users
SET is_active = $2, updated_at = NOW()
WHERE id = $1
RETURNING id, email, username, password_hash, is_active, email_verified, created_at, updated_at;

-- name: UpdateUserEmailVerified :one
UPDATE users
SET email_verified = $2, updated_at = NOW()
WHERE id = $1
RETURNING id, email, username, password_hash, is_active, email_verified, created_at, updated_at;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;