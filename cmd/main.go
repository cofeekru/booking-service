package main

import (
	"avito_tech_backend/internal/config"
	"avito_tech_backend/internal/database"
	"avito_tech_backend/internal/handlers"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/swaggo/files"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	cfg := config.MustLoad()

	connectStorage := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.DB_HOST, cfg.DB_PORT, cfg.DB_USER, cfg.DB_PASSWORD, cfg.DB_NAME, cfg.DB_SSL_MODE)
	storage, err := database.New(connectStorage)

	if err != nil {
		log.Fatalf("Failed to init storage: %s", err)
	}

	router := chi.NewRouter()

	router.Post("/dummyLogin", handlers.DummyLoginHandler())

	router.Post("/register", handlers.RegisterHandler(storage))
	router.Post("/login", handlers.LoginHandler(storage))

	router.Route("/rooms", func(r chi.Router) {
		r.Use(handlers.AuthMiddlewareHandler)

		r.Get("/list", handlers.RoomsListHandler(storage))
		r.Post("/create", handlers.RoomsCreateHandler(storage))
		r.Post("/{roomId}/schedule/create", handlers.ScheduleCreateHandler(storage))
		r.Get("/{roomId}/slots/list", handlers.SlotsListHeader(storage))
	})

	router.Route("/bookings", func(r chi.Router) {
		r.Use(handlers.AuthMiddlewareHandler)

		r.Post("/create", handlers.BookingsCreateHandler(storage))
		r.Get("/list", handlers.BookingsListHandler(storage))
		r.Get("/my", handlers.BookingsMyHandler(storage))
		r.Post("/{bookingId}/cancel", handlers.BookingsCancelHandler(storage))
	})

	router.Get("/_info", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Handle("/api.yaml", http.FileServer(http.Dir("./api")))
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/api.yaml"),
	))

	server := &http.Server{
		Addr:        cfg.HTTPServer.Host,
		Handler:     router,
		ReadTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout: cfg.HTTPServer.IdleTimeout,
	}

	slog.Info("Starting server...")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
	slog.Info("Server stopped")

}
