package main

import (
	"context"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/novaru/scallopticon/services/planet/internal/handlers"
	"github.com/novaru/scallopticon/services/planet/internal/repository"
	"github.com/novaru/scallopticon/services/planet/internal/service"
	"github.com/novaru/scallopticon/shared/db/generated"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	if err := godotenv.Load(); err != nil {
		logger.Fatal("Error loading .env file", zap.Error(err))
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		logger.Fatal("DATABASE_URL not set")
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		logger.Fatal("DB connect error", zap.Error(err))
	}
	defer pool.Close()

	q := generated.New(pool)
	repo := repository.NewPlayerRepository(q, logger)
	svc := service.NewPlayerService(repo, logger)
	handler := handlers.NewPlayerHandler(svc)

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Route("/players", func(r chi.Router) {
		r.Get("/", handler.GetPlayers)
		r.Get("/{id}", handler.GetPlayerByID)
		r.Post("/", handler.CreatePlayer)
	})

	logger.Info("Planet service running on :5000")
	if err := http.ListenAndServe(":5000", r); err != nil {
		logger.Fatal("HTTP server error", zap.Error(err))
	}
}
