package handlers

import (
	"avito_tech_backend/internal/config"
	"avito_tech_backend/internal/services"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func BookingsCreateHandler(storage config.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		userID := r.Context().Value("userID").(string)
		userRole := r.Context().Value("userRole").(string)

		var user config.User
		user.ID, _ = uuid.Parse(userID)
		user.Role = userRole

		if user.Role == "admin" {
			slog.Info("Доступ запрещён (бронирование доступно только роли user)")
			http.Error(w, string(config.FORBIDDEN), http.StatusForbidden)
			return
		}

		var booking config.Booking
		err := json.NewDecoder(r.Body).Decode(&booking)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, string(config.INTERNAL_ERROR), http.StatusInternalServerError)
			return
		}

		if !services.SlotExist(storage, booking.SlotID) {
			slog.Info("Слот не найден")
			http.Error(w, string(config.SLOT_NOT_FOUND), http.StatusNotFound)
			return
		}

		if services.SlotInPast(storage, booking.SlotID) {
			slog.Info("Неверный запрос")
			http.Error(w, string(config.INVALID_REQUEST), http.StatusBadRequest)
			return
		}

		if services.SlotBooked(storage, booking.SlotID) {
			slog.Info("Слот уже занят")
			http.Error(w, string(config.SLOT_ALREADY_BOOKED), http.StatusConflict)
			return
		}

		booking.UserID = user.ID

		err = services.CreateBooking(storage, &booking)
		if err != nil {
			http.Error(w, string(config.INTERNAL_ERROR), http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(booking)
	}
}

func BookingsListHandler(storage config.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		userID := r.Context().Value("userID").(string)
		userRole := r.Context().Value("userRole").(string)

		var user config.User
		user.ID, _ = uuid.Parse(userID)
		user.Role = userRole

		if user.Role == "user" {
			slog.Info("Доступ запрещён (только admin)")
			http.Error(w, string(config.FORBIDDEN), http.StatusForbidden)
			return
		}

		inputPage := r.URL.Query().Get("page")
		inputPageSize := r.URL.Query().Get("pageSize")

		if inputPage == "" {
			inputPage = "1"
		}

		if inputPageSize == "" {
			inputPageSize = "20"
		}

		page, _ := strconv.Atoi(inputPage)
		pageSize, _ := strconv.Atoi(inputPageSize)

		if page < 1 || pageSize < 1 || pageSize > 100 {
			http.Error(w, string(config.INVALID_REQUEST), http.StatusBadRequest)
			return
		}
		var pagination config.Pagination
		pagination.Page = page
		pagination.PageSize = pageSize

		result, err := services.BookingsList(storage, &pagination, user)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, string(config.INVALID_REQUEST), http.StatusBadRequest)
			return
		}

		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, string(config.INTERNAL_ERROR), http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(pagination)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, string(config.INTERNAL_ERROR), http.StatusInternalServerError)
			return
		}
	}
}

func BookingsMyHandler(storage config.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		userID := r.Context().Value("userID").(string)
		userRole := r.Context().Value("userRole").(string)

		var user config.User
		user.ID, _ = uuid.Parse(userID)
		user.Role = userRole

		if user.Role == "admin" {
			slog.Info("Доступ запрещён (только user)")
			http.Error(w, string(config.FORBIDDEN), http.StatusForbidden)
			return
		}

		result, err := services.BookingsMy(storage, user)
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

func BookingsCancelHandler(storage config.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		userID := r.Context().Value("userID").(string)
		userRole := r.Context().Value("userRole").(string)

		var user config.User
		user.ID, _ = uuid.Parse(userID)
		user.Role = userRole

		if user.Role == "amdin" {
			slog.Info("Роль не user")
			http.Error(w, string(config.FORBIDDEN), http.StatusForbidden)
			return
		}

		bookingIDExtract := func(url string) string {
			parts := strings.Split(url, "/")
			for i, part := range parts {
				if part == "bookings" && i+1 < len(parts) {
					return parts[i+1]
				}
			}
			return ""
		}(r.URL.Path)

		bookingID, _ := uuid.Parse(bookingIDExtract)
		if !services.BookingExist(storage, bookingID) {
			slog.Info("Бронь не найдена")
			http.Error(w, string(config.BOOKING_NOT_FOUND), http.StatusNotFound)
			return
		}

		if !services.BookingValid(storage, bookingID, user.ID) {
			slog.Info("Не своя бронь")
			http.Error(w, string(config.FORBIDDEN), http.StatusForbidden)
			return
		}

		result, err := services.CancelBooking(storage, bookingID)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, string(config.INTERNAL_ERROR), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(result)

	}
}
