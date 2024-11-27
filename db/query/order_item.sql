-- name: CreateOrderitems :one
INSERT INTO orderitems (
  item_name, item_sub_total, quantity, item_id, order
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

-- name: GetOrderitemsByOrderID :one
SELECT * FROM orderitems
WHERE order = $1 LIMIT 1;

-- name: GetOrderitemsByOrderID :many
SELECT * FROM orderitems
WHERE order = $1 
ORDER BY order;