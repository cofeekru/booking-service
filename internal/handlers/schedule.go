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

func ScheduleCreateHandler(storage config.Database) http.HandlerFunc {
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

		var schedule config.Schedule
		err := json.NewDecoder(r.Body).Decode(&schedule)

		if err != nil || !services.ValidSchedule(schedule) {
			slog.Info("Неверный запрос (в т.ч. недопустимые значения daysOfWeek)")
			slog.Info(err.Error())
			http.Error(w, string(config.INVALID_REQUEST), http.StatusBadRequest)
			return
		}

		if !services.RoomExist(storage, roomID) {
			slog.Info("Переговорка не найдена")
			http.Error(w, string(config.INVALID_REQUEST), http.StatusNotFound)
			return
		}

		if services.ScheduleExist(storage, &schedule) {
			slog.Info("Расписание для переговорки уже создано, изменение не допускается")
			http.Error(w, string(config.SCHEDULE_EXISTS), http.StatusConflict)
			return
		}

		err = services.ScheduleCreate(storage, &schedule, roomID, user)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, string(config.INTERNAL_ERROR), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(schedule)
	}
}
