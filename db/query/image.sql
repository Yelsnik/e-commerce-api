-- name: CreateImages :one
INSERT INTO images (
  image_name,
  data,
  product
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetImages :one
SELECT * FROM images
WHERE product = $1 LIMIT 1;

-- name: ListImages :many
SELECT * FROM images
ORDER BY  id
LIMIT $1
OFFSET $2;