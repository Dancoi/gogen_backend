-- name: CreateRole :one
INSERT INTO roles (name, description, permissions)
VALUES ($1, $2, $3)
RETURNING id, name, description, permissions, created_at, updated_at;

-- name: GetRoleByID :one
SELECT id, name, description, permissions, created_at, updated_at
FROM roles
WHERE id = $1;

-- name: GetRoleByName :one
SELECT id, name, description, permissions, created_at, updated_at
FROM roles
WHERE name = $1;

-- name: ListRoles :many
SELECT id, name, description, permissions, created_at, updated_at
FROM roles
ORDER BY id;

-- name: UpdateRole :one
UPDATE roles
SET name = $2, description = $3, permissions = $4, updated_at = NOW()
WHERE id = $1
RETURNING id, name, description, permissions, created_at, updated_at;

-- name: AddUserRole :exec
INSERT INTO user_roles (user_id, role_id)
VALUES ($1, $2)
ON CONFLICT (user_id, role_id) DO NOTHING;

-- name: RemoveUserRole :exec
DELETE FROM user_roles
WHERE user_id = $1 AND role_id = $2;

-- name: GetUserRoles :many
SELECT r.id, r.name, r.description, r.permissions, r.created_at, r.updated_at
FROM roles r
JOIN user_roles ur ON r.id = ur.role_id
WHERE ur.user_id = $1;

-- name: DeleteRole :exec
DELETE FROM roles
WHERE id = $1;
