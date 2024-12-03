-- name: CreateOrderitems :one
INSERT INTO orderitems (
  item_name, item_sub_total, quantity, item_id, order_id
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetOrderitems :one
SELECT * FROM orderitems
WHERE id = $1 LIMIT 1;

-- name: GetOrderitemsForUpdate :one
SELECT * FROM orderitems
WHERE id = $1
FOR NO KEY UPDATE;

-- name: GetOrderitemByOrderID :one
SELECT * FROM orderitems
WHERE order_id = $1 LIMIT 1;

-- name: GetOrderitemsByOrderID :many
SELECT * FROM orderitems
WHERE order_id = $1 
ORDER BY order_id;