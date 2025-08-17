package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"github.com/novaru/scallopticon/shared/apperrors"
	"github.com/novaru/scallopticon/shared/db/generated"
)

type PlayerRepository interface {
	GetPlayers(ctx context.Context) ([]generated.Player, error)
	GetByID(ctx context.Context, id uuid.UUID) (generated.Player, error)
	CreatePlayerWithPlanet(ctx context.Context, username, planetName string) (generated.Player, generated.Planet, error)
}

type DB interface {
	generated.DBTX
	Begin(ctx context.Context) (pgx.Tx, error)
}
type playerRepository struct {
	q      *generated.Queries
	db     DB
	logger *zap.Logger
}

func NewPlayerRepository(q *generated.Queries, logger *zap.Logger) PlayerRepository {
	return &playerRepository{
		q:      q,
		logger: logger,
	}
}

func (r *playerRepository) GetPlayers(ctx context.Context) ([]generated.Player, error) {
	players, err := r.q.ListPlayers(ctx)
	if err != nil {
		r.logger.Error("failed to list players", zap.Error(err))
		return nil, apperrors.NewInternalError("failed to retrieve players", err)
	}

	r.logger.Debug("succesfully retrieved players", zap.Int("count", len(players)))
	return players, nil
}

func (r *playerRepository) GetByID(ctx context.Context, id uuid.UUID) (generated.Player, error) {
	player, err := r.q.GetPlayerByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Debug("player not found", zap.String("player_id", id.String()))
			return generated.Player{}, apperrors.NewNotFoundError("player", "player with given ID does not exist")
		}

		r.logger.Error("faile to get player by ID",
			zap.String("player_id", id.String()),
			zap.Error(err))
		return generated.Player{}, err
	}

	return player, nil
}

func (r *playerRepository) CreatePlayerWithPlanet(ctx context.Context, username, planetName string) (generated.Player, generated.Planet, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		r.logger.Error("failed to start transaction", zap.Error(err))
		return generated.Player{}, generated.Planet{}, apperrors.NewInternalError("failed to start database transaction", err)
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				r.logger.Error("failed to rollback transaction", zap.Error(rollbackErr))
			}
		}
	}()

	qtx := r.q.WithTx(tx)

	player, err := qtx.CreatePlayer(ctx, username)
	if err != nil {
		if isDuplicateKeyError(err) {
			r.logger.Debug("player already exists", zap.String("username", username))
			return generated.Player{}, generated.Planet{},
				apperrors.NewAlreadyExistsError("player", "player with this username already exists")
		}

		r.logger.Error("failed to create player",
			zap.String("username", username),
			zap.Error(err))
		return generated.Player{}, generated.Planet{},
			apperrors.NewInternalError("failed to create player", err)
	}

	planet, err := qtx.CreatePlanet(ctx, generated.CreatePlanetParams{
		PlayerID: player.ID,
		Name:     planetName,
	})
	if err != nil {
		r.logger.Error("failed to create planet",
			zap.String("player_id", player.ID.String()),
			zap.String("planet_name", planetName),
			zap.Error(err))
		return player, generated.Planet{},
			apperrors.NewInternalError("failed to create planet", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		r.logger.Error("failed to commit transaction", zap.Error(err))
		return player, planet, apperrors.NewInternalError("failed to save player and planet", err)
	}

	r.logger.Info("successfully created player with planet",
		zap.String("player_id", player.ID.String()),
		zap.String("username", username),
		zap.String("planet_name", planetName))

	return player, planet, nil
}

func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}

	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "duplicate key") ||
		strings.Contains(errStr, "unique constraint") ||
		strings.Contains(errStr, "violates unique")
}
