-- name: CreateOrder :one
INSERT INTO orders (user_id, delivery_address, delivery_date, notes, payment_method, total_amount)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, product_id, quantity, unit_price, total_price)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetOrderByID :one
SELECT * FROM orders WHERE id = $1 AND user_id = $2;

-- name: GetOrderByIDAdmin :one
SELECT * FROM orders WHERE id = $1;

-- name: ListOrdersByUserID :many
SELECT * FROM orders
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountOrdersByUserID :one
SELECT COUNT(*) FROM orders WHERE user_id = $1;

-- name: GetOrderItems :many
SELECT
    oi.*,
    p.name      AS product_name,
    p.image_url AS product_image_url
FROM order_items oi
JOIN products p ON p.id = oi.product_id
WHERE oi.order_id = $1;

-- name: UpdateOrderStatus :one
UPDATE orders
SET status = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;
