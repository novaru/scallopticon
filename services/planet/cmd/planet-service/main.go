package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/novaru/scallopticon/services/planet/internal/handlers"
	"github.com/novaru/scallopticon/services/planet/internal/repository"
	"github.com/novaru/scallopticon/services/planet/internal/service"
	"github.com/novaru/scallopticon/shared/db/generated"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("DB connect error: %v", err)
	}
	defer pool.Close()

	q := generated.New(pool)
	repo := repository.NewPlayerRepository(q)
	svc := service.NewPlayerService(repo)
	handler := handlers.NewPlayerHandler(svc)

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/players", handler.GetPlayers)
	r.Post("/players", handler.CreatePlayer)
	r.Get("/players/{id}", handler.GetPlayerByID)

	log.Println("Planet service running on :5000")
	log.Fatal(http.ListenAndServe(":5000", r))
}
