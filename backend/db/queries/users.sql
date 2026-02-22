-- name: CreateUser :one
INSERT INTO users (first_name, last_name, phone_number, email_address)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email_address = $1 AND deleted_at IS NULL;

-- name: GetUserByPhone :one
SELECT * FROM users
WHERE phone_number = $1 AND deleted_at IS NULL;

-- name: MarkUserVerified :one
UPDATE users
SET is_verified = TRUE, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: SoftDeleteUser :exec
UPDATE users
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = $1;
