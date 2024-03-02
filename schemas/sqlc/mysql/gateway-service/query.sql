-- name: CreateUser :exec
INSERT INTO users (
    id, name, email, password
) VALUES (
    ?, ?, ?, ?
);

-- name: GetUsers :many
SELECT * FROM users;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = ?;

-- name: GetUserByName :one
SELECT * FROM users
WHERE name = ?;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = ?;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?;