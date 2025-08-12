-- name: CreatePlayer :one
INSERT INTO players (username)
VALUES ($1)
RETURNING *;

-- name: GetPlayerByID :one
SELECT * FROM players
WHERE id = $1;

-- name: ListPlayers :many
SELECT * FROM players
ORDER BY created_at DESC;
