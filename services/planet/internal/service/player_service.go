package service

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/novaru/scallopticon/services/planet/internal/repository"
	"github.com/novaru/scallopticon/shared/db/generated"
)

type PlayerResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

type PlanetResponse struct {
	ID       uuid.UUID `json:"id"`
	PlayerID uuid.UUID `json:"player_id"`
	Name     string    `json:"name"`
}

type CreatePlayerResponse struct {
	Player PlayerResponse `json:"player"`
	Planet PlanetResponse `json:"planet"`
}

type PlayerService interface {
	GetAllPlayers(ctx context.Context) ([]PlayerResponse, error)
	GetPlayerByID(ctx context.Context, id uuid.UUID) (PlayerResponse, error)
	CreatePlayerWithPlanet(ctx context.Context, username, planetName string) (CreatePlayerResponse, error)
}

type playerService struct {
	repo   repository.PlayerRepository
	logger *zap.Logger
}

func NewPlayerService(repo repository.PlayerRepository, logger *zap.Logger) PlayerService {
	return &playerService{
		repo:   repo,
		logger: logger,
	}
}

func (s *playerService) CreatePlayerWithPlanet(ctx context.Context, username, planetName string) (CreatePlayerResponse, error) {
	s.logger.Debug("creating player with planet",
		zap.String("username", username),
		zap.String("planet_name", planetName))

	// TODO: validate and sanitize inputs

	normalizedUsername := s.normalizeUsername(username)
	normalizedPlanetName := strings.TrimSpace(planetName)

	player, planet, err := s.repo.CreatePlayerWithPlanet(ctx, normalizedUsername, normalizedPlanetName)
	if err != nil {
		return CreatePlayerResponse{}, err
	}

	response := CreatePlayerResponse{
		Player: s.convertPlayerToResponse(player),
		Planet: s.convertPlanetToResponse(planet),
	}

	s.logger.Info("successfully created player with planet",
		zap.String("player_id", player.ID.String()),
		zap.String("username", normalizedUsername),
		zap.String("planet_name", normalizedPlanetName))

	return response, nil
}

func (s *playerService) GetAllPlayers(ctx context.Context) ([]PlayerResponse, error) {
	s.logger.Debug("retrieving all players")

	players, err := s.repo.GetPlayers(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]PlayerResponse, len(players))
	for i, player := range players {
		responses[i] = s.convertPlayerToResponse(player)
	}

	s.logger.Debug("successfully retrieved all players", zap.Int("count", len(responses)))
	return responses, nil
}

func (s *playerService) GetPlayerByID(ctx context.Context, id uuid.UUID) (PlayerResponse, error) {
	s.logger.Debug("retrieving player by ID", zap.String("player_id", id.String()))

	player, err := s.repo.GetByID(ctx, id)
	if err != nil {
		// Error already logged and wrapped in repository
		return PlayerResponse{}, err
	}

	response := s.convertPlayerToResponse(player)
	s.logger.Debug("successfully retrieved player", zap.String("player_id", id.String()))

	return response, nil
}

func (s *playerService) normalizeUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

// Convert generated models to domain response models
func (s *playerService) convertPlayerToResponse(player generated.Player) PlayerResponse {
	return PlayerResponse{
		ID:        player.ID,
		Username:  player.Username,
		CreatedAt: player.CreatedAt.Time, // Assuming CreatedAt is sql.NullTime
	}
}

func (s *playerService) convertPlanetToResponse(planet generated.Planet) PlanetResponse {
	return PlanetResponse{
		ID:       planet.ID,
		PlayerID: planet.PlayerID,
		Name:     planet.Name,
	}
}
