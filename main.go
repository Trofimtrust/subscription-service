package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "subscription-service/docs"
	"subscription-service/internal/config"
	"subscription-service/internal/handler"
	"subscription-service/internal/repository"
)

// @title Subscription Service API
// @version 1.0
// @description REST API для управления подписками.
// @host localhost:8080
// @BasePath /

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	conn, err := config.ConnectDB()
	if err != nil {
		logger.Error("Database connection failed", "error", err)
		return
	}
	defer conn.Close(context.Background())

	repo := repository.New(conn)
	h := handler.New(repo)

	r := chi.NewRouter()

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Subscription Service is running!"))
	})

	r.Post("/subscriptions", h.CreateSubscription)
	r.Get("/subscriptions", h.GetSubscriptions)
	r.Get("/subscriptions/cost", h.CalculateCost)
	r.Get("/subscriptions/{id}", h.GetSubscriptionByID)
	r.Put("/subscriptions/{id}", h.UpdateSubscription)
	r.Delete("/subscriptions/{id}", h.DeleteSubscription)

	logger.Info("Server started", "port", 8080)

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		logger.Error("Server stopped", "error", err)
	}
}
