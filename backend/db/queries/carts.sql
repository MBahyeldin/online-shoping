-- name: GetOrCreateCart :one
INSERT INTO carts (user_id)
VALUES ($1)
ON CONFLICT (user_id) DO UPDATE SET updated_at = NOW()
RETURNING *;

-- name: GetCartByUserID :one
SELECT * FROM carts WHERE user_id = $1;

-- name: GetCartItems :many
SELECT
    ci.id,
    ci.cart_id,
    ci.product_id,
    ci.quantity,
    ci.created_at,
    ci.updated_at,
    p.name        AS product_name,
    p.price       AS product_price,
    p.image_url   AS product_image_url,
    p.stock_quantity AS product_stock
FROM cart_items ci
JOIN products p ON p.id = ci.product_id
WHERE ci.cart_id = $1
ORDER BY ci.created_at ASC;

-- name: UpsertCartItem :one
INSERT INTO cart_items (cart_id, product_id, quantity)
VALUES ($1, $2, $3)
ON CONFLICT (cart_id, product_id)
DO UPDATE SET quantity = $3, updated_at = NOW()
RETURNING *;

-- name: UpdateCartItemQuantity :one
UPDATE cart_items
SET quantity = $2, updated_at = NOW()
WHERE id = $1 AND cart_id = $3
RETURNING *;

-- name: DeleteCartItem :exec
DELETE FROM cart_items WHERE id = $1 AND cart_id = $2;

-- name: ClearCart :exec
DELETE FROM cart_items WHERE cart_id = $1;

-- name: GetCartItemByID :one
SELECT ci.*, p.name AS product_name, p.price AS product_price
FROM cart_items ci
JOIN products p ON p.id = ci.product_id
WHERE ci.id = $1;
