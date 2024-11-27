-- name: CreateCarts :one
INSERT INTO carts (
 user_id, total_price 
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetCartsByUserID :one
SELECT * FROM carts
WHERE user_id = $1 LIMIT 1;

-- name: GetCarts :one
SELECT * FROM carts
WHERE id = $1 LIMIT 1;

-- name: GetCartsForUpdate :one
SELECT * FROM carts
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: UpdateCarts :one
UPDATE carts
  set total_price = $2
WHERE id = $1
RETURNING *;