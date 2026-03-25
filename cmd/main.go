package main

import (
	"avito_tech_backend/internal/config"
	"avito_tech_backend/internal/database"
	"avito_tech_backend/internal/handlers"
	"log"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.MustLoad("./config/config.yaml")
	storage, err := database.New(cfg.ConnectStorage)

	if err != nil {
		log.Fatalf("Failed to init storage: %s", err)
	}

	router := chi.NewRouter()

	router.Post("/dummyLogin", handlers.DummyLoginHandler())

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
		r.Get("/{bookingId}/cancel", handlers.BookingsCancelHandler(storage))
	})

	router.Get("/_info", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	server := &http.Server{
		Addr:        cfg.Address,
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
