package types

import "time"

type SimulationRequest struct {
	PlanetID string `json:"planet_id"`
	WaveID   string `json:"wave_id"`
}

type SimulationResult struct {
	DamageTaken      int       `json:"damage_taken"`
	ShieldsRemaining int       `json:"shields_remaining"`
	HPRemaining      int       `json:"hp_remaining"`
	AliensDestroyed  int       `json:"aliens_destroyed"`
	Loot             Resources `json:"loot"`
	Events           []string  `json:"events,omitempty"`
	Timestamp        time.Time `json:"timestamp"`
}
