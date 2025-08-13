package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/novaru/scallopticon/services/planet/internal/service"
	"github.com/novaru/scallopticon/shared/apperrors"
	"github.com/novaru/scallopticon/shared/response"
)

type PlayerHandler struct {
	service service.PlayerService
}

func NewPlayerHandler(s service.PlayerService) *PlayerHandler {
	return &PlayerHandler{service: s}
}

type CreatePlayerRequest struct {
	Username   string `json:"username"`
	PlanetName string `json:"planet_name"`
}

func (r *CreatePlayerRequest) Validate() error {
	if strings.TrimSpace(r.Username) == "" {
		return apperrors.NewInvalidInputError("username is required", nil)
	}
	if len(r.Username) < 3 {
		return apperrors.NewInvalidInputError("username must be at least 3 characters", nil)
	}
	if strings.TrimSpace(r.PlanetName) == "" {
		return apperrors.NewInvalidInputError("planet name is required", nil)
	}
	return nil
}

func (h *PlayerHandler) GetPlayers(w http.ResponseWriter, r *http.Request) {
	players, err := h.service.GetAllPlayers(r.Context())
	if err != nil {
		response.WriteError(w, err)
		return
	}

	response.WriteSuccess(w, players)
}

func (h *PlayerHandler) GetPlayerByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	playerID, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "invalid UUID", http.StatusBadRequest)
		return
	}

	player, err := h.service.GetPlayerByID(r.Context(), playerID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Player not found", http.StatusNotFound)
			return
		}
		log.Println("Error fetching player:", err)
		http.Error(w, "Failed to fetch player", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(player)
}

func (h *PlayerHandler) CreatePlayer(w http.ResponseWriter, r *http.Request) {
	var req CreatePlayerRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		response.WriteError(w, apperrors.NewInvalidInputError("invalid JSON format", err))
		return
	}

	if err := req.Validate(); err != nil {
		response.WriteError(w, err)
		return
	}

	result, err := h.service.CreatePlayerWithPlanet(r.Context(), req.Username, req.PlanetName)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	response.WriteCreated(w, result)
}
