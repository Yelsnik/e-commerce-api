-- name: CreateCartitems :one
INSERT INTO cartitems (
  cart, product, quantity, price, sub_total
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetCartitems :one
SELECT * FROM cartitems
WHERE id = $1 LIMIT 1;

-- name: ListCartitems :many
SELECT * FROM cartitems
ORDER BY  id
LIMIT $1
OFFSET $2;

-- name: UpdateCartitems :one
UPDATE cartitems
  set quantity = $2,
  price = $3,
  sub_total = $4
WHERE id = $1
RETURNING *;