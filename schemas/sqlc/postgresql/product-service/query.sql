-- name: CreateProduct :exec
INSERT INTO products (
    id, user_id, name, price, description, stock
) VALUES (
    $1, $2, $3, $4, $5, $6
);

-- name: GetProducts :many
SELECT * FROM products;

-- name: GetProductsByUserId :many
SELECT * FROM products
WHERE user_id = $1;

-- name: GetProductById :one
SELECT * FROM products
WHERE id = $1;

-- name: UpdateProduct :exec
UPDATE products
    set price = $2,
    description = $3,
    stock = $4
WHERE id = $1;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1;
