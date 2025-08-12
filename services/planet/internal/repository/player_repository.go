package repository

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/novaru/scallopticon/shared/db/generated"
)

type PlayerRepository struct {
	q *generated.Queries
}

func NewPlayerRepository(q *generated.Queries) *PlayerRepository {
	return &PlayerRepository{q: q}
}

func (r *PlayerRepository) GetPlayers(ctx context.Context) ([]generated.Player, error) {
	players, err := r.q.ListPlayers(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return players, nil
}

func (r *PlayerRepository) GetByID(ctx context.Context, id uuid.UUID) (generated.Player, error) {
	player, err := r.q.GetPlayerByID(ctx, id)
	if err != nil {
		return generated.Player{}, err
	}

	return player, nil
}

func (r *PlayerRepository) CreatePlayerWithPlanet(ctx context.Context, username, planetName string) (generated.Player, generated.Planet, error) {
	player, err := r.q.CreatePlayer(ctx, username)
	if err != nil {
		return generated.Player{}, generated.Planet{}, err
	}

	planet, err := r.q.CreatePlanet(ctx, generated.CreatePlanetParams{
		PlayerID: player.ID,
		Name:     planetName,
	})
	if err != nil {
		return player, generated.Planet{}, err
	}

	return player, planet, nil
}
