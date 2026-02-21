-- name: CreateOrder :one
INSERT INTO orders (coupon_code, total, discounts)
VALUES ($1, $2, $3)
RETURNING id, coupon_code, total, discounts, created_at;

-- name: CreateOrderItem :exec
INSERT INTO order_items (order_id, product_id, quantity)
VALUES ($1, $2, $3);

-- name: GetOrder :one
SELECT id, coupon_code, total, discounts, created_at
FROM orders
WHERE id = $1;

-- name: GetOrderItems :many
SELECT product_id, quantity
FROM order_items
WHERE order_id = $1;
