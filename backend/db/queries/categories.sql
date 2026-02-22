-- name: CreateCategory :one
INSERT INTO categories (name, slug)
VALUES ($1, $2)
RETURNING *;

-- name: ListCategories :many
SELECT * FROM categories
ORDER BY name ASC;

-- name: GetCategoryByID :one
SELECT * FROM categories
WHERE id = $1;

-- name: GetCategoryBySlug :one
SELECT * FROM categories
WHERE slug = $1;
