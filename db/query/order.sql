-- name: CreateOrders :one
INSERT INTO orders (
  user_name, user_id, total_price, delivery_address, country, status
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetOrdersByUserID :one
SELECT * FROM orders
WHERE user_id = $1 LIMIT 1;

-- name: GetOrdersByID :one
SELECT * FROM orders
WHERE id = $1 LIMIT 1;

-- name: GetOrdersForUpdate :one
SELECT * FROM orders
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: UpdateOrders :one
UPDATE orders
  set total_price = $2,
  delivery_address = $3,
  country = $4,
  status = $5
WHERE id = $1
RETURNING *;