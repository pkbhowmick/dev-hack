// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: query.sql

package db

import (
	"context"
)

const createProduct = `-- name: CreateProduct :exec
INSERT INTO products (
    id, user_id, name, price, description, stock
) VALUES (
    $1, $2, $3, $4, $5, $6
)
`

type CreateProductParams struct {
	ID          string `json:"id"`
	UserID      string `json:"userId"`
	Name        string `json:"name"`
	Price       int32  `json:"price"`
	Description string `json:"description"`
	Stock       int32  `json:"stock"`
}

func (q *Queries) CreateProduct(ctx context.Context, arg CreateProductParams) error {
	_, err := q.db.ExecContext(ctx, createProduct,
		arg.ID,
		arg.UserID,
		arg.Name,
		arg.Price,
		arg.Description,
		arg.Stock,
	)
	return err
}

const deleteProduct = `-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1
`

func (q *Queries) DeleteProduct(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, deleteProduct, id)
	return err
}

const getProductById = `-- name: GetProductById :one
SELECT id, user_id, name, price, description, stock, created_at, updated_at FROM products
WHERE id = $1
`

func (q *Queries) GetProductById(ctx context.Context, id string) (Product, error) {
	row := q.db.QueryRowContext(ctx, getProductById, id)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.Price,
		&i.Description,
		&i.Stock,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getProducts = `-- name: GetProducts :many
SELECT id, user_id, name, price, description, stock, created_at, updated_at FROM products
`

func (q *Queries) GetProducts(ctx context.Context) ([]Product, error) {
	rows, err := q.db.QueryContext(ctx, getProducts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Product
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Name,
			&i.Price,
			&i.Description,
			&i.Stock,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const getProductsByUserId = `-- name: GetProductsByUserId :many
SELECT id, user_id, name, price, description, stock, created_at, updated_at FROM products
WHERE user_id = $1
`

func (q *Queries) GetProductsByUserId(ctx context.Context, userID string) ([]Product, error) {
	rows, err := q.db.QueryContext(ctx, getProductsByUserId, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Product
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Name,
			&i.Price,
			&i.Description,
			&i.Stock,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const updateProduct = `-- name: UpdateProduct :exec
UPDATE products
    set price = $2,
    description = $3,
    stock = $4
WHERE id = $1
`

type UpdateProductParams struct {
	ID          string `json:"id"`
	Price       int32  `json:"price"`
	Description string `json:"description"`
	Stock       int32  `json:"stock"`
}

func (q *Queries) UpdateProduct(ctx context.Context, arg UpdateProductParams) error {
	_, err := q.db.ExecContext(ctx, updateProduct,
		arg.ID,
		arg.Price,
		arg.Description,
		arg.Stock,
	)
	return err
}
