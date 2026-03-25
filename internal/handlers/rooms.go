package handlers

import (
	"avito_tech_backend/internal/config"
	"avito_tech_backend/internal/services"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

func RoomsListHandler(storage config.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		result, err := services.RoomsList(storage)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, string(config.INTERNAL_ERROR), http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, string(config.INTERNAL_ERROR), http.StatusInternalServerError)
			return
		}
	}
}

func RoomsCreateHandler(storage config.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		userID := r.Context().Value("userID").(string)
		userRole := r.Context().Value("userRole").(string)

		var user config.User
		user.ID, _ = uuid.Parse(userID)
		user.Role = userRole

		if user.Role == "user" {
			slog.Info("Доступ запрещён (требуется роль admin)")
			http.Error(w, string(config.FORBIDDEN), http.StatusForbidden)
			return
		}

		var room config.Room
		err := json.NewDecoder(r.Body).Decode(&room)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, string(config.INTERNAL_ERROR), http.StatusInternalServerError)
			return
		}

		err = services.RoomsCreate(storage, user, &room)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, string(config.INTERNAL_ERROR), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(room)
	}
}
