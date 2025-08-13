package repository

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/novaru/scallopticon/shared/db/generated"
)

type PlayerRepository interface {
	GetPlayers(ctx context.Context) ([]generated.Player, error)
	GetByID(ctx context.Context, id uuid.UUID) (generated.Player, error)
	CreatePlayerWithPlanet(ctx context.Context, username, planetName string) (generated.Player, generated.Planet, error)
}

type playerRepository struct {
	q *generated.Queries
}

func NewPlayerRepository(q *generated.Queries) PlayerRepository {
	return &playerRepository{q: q}
}

func (r *playerRepository) GetPlayers(ctx context.Context) ([]generated.Player, error) {
	players, err := r.q.ListPlayers(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return players, nil
}

func (r *playerRepository) GetByID(ctx context.Context, id uuid.UUID) (generated.Player, error) {
	player, err := r.q.GetPlayerByID(ctx, id)
	if err != nil {
		return generated.Player{}, err
	}

	return player, nil
}

func (r *playerRepository) CreatePlayerWithPlanet(ctx context.Context, username, planetName string) (generated.Player, generated.Planet, error) {
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
