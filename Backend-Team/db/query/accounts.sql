-- name: GetAccount :one
SELECT * FROM accounts 
WHERE username = $1 LIMIT 1;

-- name: ListAccounts :many
SELECT * FROM accounts
ORDER BY id LIMIT $1 OFFSET $2;

-- name: UpdateAccount :one
UPDATE accounts
SET balance = $2
WHERE username = $1
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE username = $1;

-- name: CreateAccount :one
INSERT INTO accounts (username, balance, currency)
VALUES ($1, $2, $3)
RETURNING *;
