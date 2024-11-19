-- name: CreateCarts :one
INSERT INTO carts (
 user_id, total_price 
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetCarts :one
SELECT * FROM carts
WHERE id = $1 LIMIT 1;

-- name: UpdateCarts :one
UPDATE carts
  set total_price = $2
WHERE id = $1
RETURNING *;