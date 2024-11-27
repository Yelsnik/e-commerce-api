-- name: CreateProducts :one
INSERT INTO products (
  category,
  product_name,
  description,
  brand,
  count_in_stock,
  price,
  currency,
  rating,
  is_featured,
  user_id
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- name: GetProducts :one
SELECT * FROM products
WHERE id = $1 LIMIT 1;

-- name: ListProducts :many
SELECT * FROM products
ORDER BY  id
LIMIT $1
OFFSET $2;

-- name: UpdateProducts :one
UPDATE products
  set category = $2,
  product_name = $3,
  description = $4,
  brand = $5,
  count_in_stock = $6,
  price = $7,
  rating = $8,
  is_featured = $9
WHERE id = $1
RETURNING *;

-- name: DeleteProducts :exec
DELETE FROM products
WHERE id = $1;