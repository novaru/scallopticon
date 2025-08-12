package types

import (
	"time"
)

type Resources struct {
	Minerals  int `json:"minerals" db:"minerals"`
	Energy    int `json:"energy" db:"energy"`
	TechParts int `json:"tech_parts" db:"tech_parts"`
}

type Planet struct {
	ID          string          `json:"id" db:"id"`
	Name        string          `json:"name" db:"name"`
	HP          int             `json:"hp" db:"hp"`
	Shields     int             `json:"shields" db:"shields"`
	Resources   Resources       `json:"resources" db:"resources"` // JSONB
	Defenses    []DefenseSystem `json:"defenses,omitempty" db:"-"`
	LastUpdated time.Time       `json:"last_updated" db:"last_updated"`
}

type DefenseSystem struct {
	ID          string    `json:"id" db:"id"`
	PlanetID    string    `json:"planet_id" db:"planet_id"`
	Name        string    `json:"name" db:"name"`
	Damage      int       `json:"damage" db:"damage"`
	Range       int       `json:"range" db:"range"`
	FireRate    float64   `json:"fire_rate" db:"fire_rate"`
	Level       int       `json:"level" db:"level"`
	UpgradeCost Resources `json:"upgrade_cost" db:"upgrade_cost"` // JSONB
}
