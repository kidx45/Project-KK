-- name: CreateUser :one
INSERT INTO users (username, hashed_password, email, full_name, phone_number)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id LIMIT $1 OFFSET $2;

-- name: UpdateUserOTP :one
UPDATE users
SET is_email_verified = $2, is_phone_number_verified = $3
WHERE username = $1
RETURNING *;

-- name: UpdateUserPassword :one
UPDATE users
SET hashed_password = $2, password_changed_at = $3
WHERE username = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE username = $1;