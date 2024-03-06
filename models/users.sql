-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUsers :many
SELECT * FROM users ORDER BY id;

-- name: IsEmailTaken :one
SELECT 1 FROM users WHERE email = $1;

-- name: CreateUser :one
INSERT INTO users (
  name,
  email,
  password,
  verification_token,
  avatar
) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: VerifyUser :one
UPDATE users SET verified = true, verification_token = NULL WHERE id = $1 RETURNING id;

-- name: DeactivateUser :one
UPDATE users SET active = false WHERE id = $1 RETURNING id;
