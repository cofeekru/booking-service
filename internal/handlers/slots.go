package handlers

import (
	"avito_tech_backend/internal/config"
	"avito_tech_backend/internal/services"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func SlotsListHeader(storage config.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		roomIDExtract := func(url string) string {
			parts := strings.Split(url, "/")
			for i, part := range parts {
				if part == "rooms" && i+1 < len(parts) {
					return parts[i+1]
				}
			}
			return ""
		}(r.URL.Path)

		roomID, _ := uuid.Parse(roomIDExtract)
		dateRequest := r.URL.Query().Get("date")

		if dateRequest == "" {
			slog.Error("Empty date")
			http.Error(w, string(config.INVALID_REQUEST), http.StatusBadRequest)
			return
		}

		if !services.RoomExist(storage, roomID) {
			slog.Error("Room doesn't exist")
			http.Error(w, string(config.INVALID_REQUEST), http.StatusNotFound)
			return
		}

		result, err := services.SlotsList(storage, roomID, dateRequest)
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
