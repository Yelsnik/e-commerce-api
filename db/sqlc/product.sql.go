// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: product.sql

package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createProducts = `-- name: CreateProducts :one
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
) RETURNING id, category, product_name, description, brand, count_in_stock, price, currency, rating, is_featured, user_id, created_at
`

type CreateProductsParams struct {
	Category     string         `json:"category"`
	ProductName  string         `json:"product_name"`
	Description  string         `json:"description"`
	Brand        sql.NullString `json:"brand"`
	CountInStock int64          `json:"count_in_stock"`
	Price        float64        `json:"price"`
	Currency     string         `json:"currency"`
	Rating       sql.NullInt64  `json:"rating"`
	IsFeatured   sql.NullBool   `json:"is_featured"`
	UserID       uuid.UUID      `json:"user_id"`
}

func (q *Queries) CreateProducts(ctx context.Context, arg CreateProductsParams) (Product, error) {
	row := q.db.QueryRowContext(ctx, createProducts,
		arg.Category,
		arg.ProductName,
		arg.Description,
		arg.Brand,
		arg.CountInStock,
		arg.Price,
		arg.Currency,
		arg.Rating,
		arg.IsFeatured,
		arg.UserID,
	)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Category,
		&i.ProductName,
		&i.Description,
		&i.Brand,
		&i.CountInStock,
		&i.Price,
		&i.Currency,
		&i.Rating,
		&i.IsFeatured,
		&i.UserID,
		&i.CreatedAt,
	)
	return i, err
}

const deleteProducts = `-- name: DeleteProducts :exec
DELETE FROM products
WHERE id = $1
`

func (q *Queries) DeleteProducts(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteProducts, id)
	return err
}

const getProductForUpdate = `-- name: GetProductForUpdate :one
SELECT id, category, product_name, description, brand, count_in_stock, price, currency, rating, is_featured, user_id, created_at FROM products
WHERE id = $1
FOR NO KEY UPDATE
`

func (q *Queries) GetProductForUpdate(ctx context.Context, id uuid.UUID) (Product, error) {
	row := q.db.QueryRowContext(ctx, getProductForUpdate, id)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Category,
		&i.ProductName,
		&i.Description,
		&i.Brand,
		&i.CountInStock,
		&i.Price,
		&i.Currency,
		&i.Rating,
		&i.IsFeatured,
		&i.UserID,
		&i.CreatedAt,
	)
	return i, err
}

const getProducts = `-- name: GetProducts :one
SELECT id, category, product_name, description, brand, count_in_stock, price, currency, rating, is_featured, user_id, created_at FROM products
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetProducts(ctx context.Context, id uuid.UUID) (Product, error) {
	row := q.db.QueryRowContext(ctx, getProducts, id)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Category,
		&i.ProductName,
		&i.Description,
		&i.Brand,
		&i.CountInStock,
		&i.Price,
		&i.Currency,
		&i.Rating,
		&i.IsFeatured,
		&i.UserID,
		&i.CreatedAt,
	)
	return i, err
}

const listProducts = `-- name: ListProducts :many
SELECT id, category, product_name, description, brand, count_in_stock, price, currency, rating, is_featured, user_id, created_at FROM products
ORDER BY  id
LIMIT $1
OFFSET $2
`

type ListProductsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListProducts(ctx context.Context, arg ListProductsParams) ([]Product, error) {
	rows, err := q.db.QueryContext(ctx, listProducts, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Product
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ID,
			&i.Category,
			&i.ProductName,
			&i.Description,
			&i.Brand,
			&i.CountInStock,
			&i.Price,
			&i.Currency,
			&i.Rating,
			&i.IsFeatured,
			&i.UserID,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateProducts = `-- name: UpdateProducts :one
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
RETURNING id, category, product_name, description, brand, count_in_stock, price, currency, rating, is_featured, user_id, created_at
`

type UpdateProductsParams struct {
	ID           uuid.UUID      `json:"id"`
	Category     string         `json:"category"`
	ProductName  string         `json:"product_name"`
	Description  string         `json:"description"`
	Brand        sql.NullString `json:"brand"`
	CountInStock int64          `json:"count_in_stock"`
	Price        float64        `json:"price"`
	Rating       sql.NullInt64  `json:"rating"`
	IsFeatured   sql.NullBool   `json:"is_featured"`
}

func (q *Queries) UpdateProducts(ctx context.Context, arg UpdateProductsParams) (Product, error) {
	row := q.db.QueryRowContext(ctx, updateProducts,
		arg.ID,
		arg.Category,
		arg.ProductName,
		arg.Description,
		arg.Brand,
		arg.CountInStock,
		arg.Price,
		arg.Rating,
		arg.IsFeatured,
	)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Category,
		&i.ProductName,
		&i.Description,
		&i.Brand,
		&i.CountInStock,
		&i.Price,
		&i.Currency,
		&i.Rating,
		&i.IsFeatured,
		&i.UserID,
		&i.CreatedAt,
	)
	return i, err
}
