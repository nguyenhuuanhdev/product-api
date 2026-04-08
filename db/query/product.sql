-- name: CreateProduct :one
INSERT INTO products (name, price, image)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetProducts :many
SELECT * FROM products;

-- name: GetProductByID :one
SELECT * FROM products WHERE id = $1;

-- name: UpdateProduct :one
UPDATE products
SET name = $1, price = $2, image = $3
WHERE id = $4
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1;

-- name: SearchProducts :many
SELECT id, name, price, image FROM products
WHERE LOWER(name) LIKE '%' || $1::text || '%';

-- name: SortProductsByPriceAsc :many
SELECT * FROM products
ORDER BY price ASC;

-- name: SortProductsByPriceDesc :many
SELECT * FROM products
ORDER BY price DESC;