-- name: CreateProduct :one
INSERT INTO products (category_id, name, description, price, image_url, stock_quantity)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetProductByID :one
SELECT p.*, c.name AS category_name, c.slug AS category_slug
FROM products p
LEFT JOIN categories c ON c.id = p.category_id
WHERE p.id = $1 AND p.deleted_at IS NULL;

-- name: ListProducts :many
SELECT p.*, c.name AS category_name, c.slug AS category_slug
FROM products p
LEFT JOIN categories c ON c.id = p.category_id
WHERE p.deleted_at IS NULL
  AND p.is_active = TRUE
  AND ($1::uuid IS NULL OR p.category_id = $1)
ORDER BY
  CASE WHEN $2::text = 'price_asc'  THEN p.price END ASC,
  CASE WHEN $2::text = 'price_desc' THEN p.price END DESC,
  p.created_at DESC
LIMIT $3 OFFSET $4;

-- name: CountProducts :one
SELECT COUNT(*) FROM products
WHERE deleted_at IS NULL
  AND is_active = TRUE
  AND ($1::uuid IS NULL OR category_id = $1);

-- name: DeductProductStock :exec
UPDATE products
SET stock_quantity = stock_quantity - $2, updated_at = NOW()
WHERE id = $1 AND stock_quantity >= $2;

-- name: GetProductsForOrder :many
SELECT id, name, price, stock_quantity
FROM products
WHERE id = ANY($1::uuid[])
  AND deleted_at IS NULL
  AND is_active = TRUE;

-- name: UpdateProduct :one
UPDATE products
SET name = $2, description = $3, price = $4, image_url = $5,
    stock_quantity = $6, category_id = $7, is_active = $8, updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteProduct :exec
UPDATE products
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = $1;
