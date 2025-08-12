package service

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/novaru/scallopticon/services/planet/internal/repository"
	"github.com/novaru/scallopticon/shared/db/generated"
)

type PlayerService interface {
	GetAllPlayers(ctx context.Context) ([]generated.Player, error)
	GetPlayerByID(ctx context.Context, id uuid.UUID) (generated.Player, error)
	CreatePlayerWithPlanet(ctx context.Context, username, planetName string) (generated.Player, generated.Planet, error)
}

type playerService struct {
	repo *repository.PlayerRepository
}

func NewPlayerService(repo *repository.PlayerRepository) PlayerService {
	return &playerService{repo: repo}
}

func (s *playerService) CreatePlayerWithPlanet(ctx context.Context, username, planetName string) (generated.Player, generated.Planet, error) {
	username = strings.ToLower(username)
	return s.repo.CreatePlayerWithPlanet(ctx, username, planetName)
}

func (s *playerService) GetAllPlayers(ctx context.Context) ([]generated.Player, error) {
	return s.repo.GetPlayers(ctx)
}

func (s *playerService) GetPlayerByID(ctx context.Context, id uuid.UUID) (generated.Player, error) {
	return s.repo.GetByID(ctx, id)
}
