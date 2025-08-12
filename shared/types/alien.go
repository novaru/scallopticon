package types

import "time"

type AlienTemplate struct {
	ID           string             `json:"id" db:"id"`
	Name         string             `json:"name" db:"name"`
	HP           int                `json:"hp" db:"hp"`
	Damage       int                `json:"damage" db:"damage"`
	Speed        float64            `json:"speed" db:"speed"`
	BehaviorType string             `json:"behavior_type" db:"behavior_type"` // used to instantiate behavior
	Resistances  map[string]float64 `json:"resistances" db:"resistances"`     // JSONB
	LootDrop     Resources          `json:"loot_drop" db:"loot_drop"`         // JSONB
}

type Wave struct {
	ID         string      `json:"id" db:"id"`
	Difficulty int         `json:"difficulty" db:"difficulty"`
	Aliens     []WaveSpawn `json:"aliens" db:"-"`
	CreatedAt  time.Time   `json:"created_at" db:"created_at"`
}

type WaveSpawn struct {
	AlienID string `json:"alien_id" db:"alien_id"`
	Count   int    `json:"count" db:"count"`
}
