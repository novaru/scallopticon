package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/novaru/scallopticon/services/planet/internal/service"
)

type PlayerHandler struct {
	service service.PlayerService
}

func NewPlayerHandler(s service.PlayerService) *PlayerHandler {
	return &PlayerHandler{service: s}
}

func (h *PlayerHandler) GetPlayers(w http.ResponseWriter, r *http.Request) {
	players, err := h.service.GetAllPlayers(r.Context())
	if err != nil {
		http.Error(w, "Failed to queries all players", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(players)
}

func (h *PlayerHandler) GetPlayerByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	playerID, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "invalid UUID", http.StatusBadRequest)
		return
	}

	player, err := h.service.GetPlayerByID(r.Context(), playerID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(player)
}

func (h *PlayerHandler) CreatePlayer(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username   string `json:"username"`
		PlanetName string `json:"planet_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	player, planet, err := h.service.CreatePlayerWithPlanet(r.Context(), req.Username, req.PlanetName)
	if err != nil {
		log.Println("Error creating player with planet:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"player": player,
		"planet": planet,
	})
}
