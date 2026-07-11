-- name: CreateSession :one
INSERT INTO sessions (id, username, user_agent, client_ip, refresh_token, is_blocked, expires_at) 
VALUES ($1, $2, $3, $4, $5, $6, $7) 
RETURNING *;

-- name: GetSessionByUsername :one
SELECT 
* FROM sessions WHERE username = $1 LIMIT 1;

-- name: GetSessionById :one
SELECT 
* FROM sessions WHERE id = $1 LIMIT 1;