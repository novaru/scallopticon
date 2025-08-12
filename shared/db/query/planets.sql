-- name: CreatePlanet :one
INSERT INTO planets (player_id, name)
VALUES ($1, $2)
RETURNING *;

-- name: GetPlanetByPlayerID :one
SELECT * FROM planets
WHERE player_id = $1;

-- name: UpdatePlanetState :exec
UPDATE planets
SET resources = $2,
    defense_level = $3,
    current_wave = $4,
    health = $5,
    updated_at = now()
WHERE id = $1;

-- name: DeletePlanet :exec
DELETE FROM planets
WHERE id = $1;
