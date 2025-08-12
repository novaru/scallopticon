-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE players (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username    TEXT NOT NULL UNIQUE,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE TABLE planets (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id       UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    name            TEXT NOT NULL,
    resources       INT DEFAULT 0,
    defense_level   INT DEFAULT 1,
    current_wave    INT DEFAULT 0,
    health          INT DEFAULT 100,
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT now(),
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- index for fast lookups by player
CREATE INDEX idx_planets_player_id ON planets(player_id);


-- +goose Down
ALTER TABLE planets DROP CONSTRAINT planets_player_id_fkey;
DROP TABLE IF EXISTS planets;
DROP TABLE IF EXISTS players;

