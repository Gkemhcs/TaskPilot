-- name: CreateUser :one 
INSERT INTO users(name,hashed_password) VALUES ($1,$2) RETURNING *;


-- name: GetUserById :one

SELECT * FROM users WHERE id=$1;

-- name: GetUserByName :one
SELECT * FROM users WHERE name=$1;


-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY id;