package services

import (
	"avito_tech_backend/internal/config"
	"time"

	"github.com/google/uuid"
)

func CreateBooking(storage config.Database, booking *config.Booking) error {
	booking.ID = uuid.New()
	booking.Status = "active"
	booking.CreatedAt = time.Now().Format(time.RFC3339)
	err := storage.CreateBooking(booking)

	return err
}

func BookingsList(storage config.Database, pagination *config.Pagination, user config.User) ([]config.Booking, error) {
	result, err := storage.GetAdminBookingList(pagination, user)
	return result, err
}

func BookingsMy(storage config.Database, user config.User) ([]config.Booking, error) {
	result, err := storage.GetUserBookingList(user)
	return result, err
}

func BookingExist(storage config.Database, bookingID uuid.UUID) bool {
	_, err := storage.GetBookingByBookingID(bookingID)
	return err == nil
}

func BookingValid(storage config.Database, bookingID uuid.UUID, userID uuid.UUID) bool {
	result, err := storage.GetBookingByBookingID(bookingID)
	if err != nil {
		return false
	}
	return result.UserID == userID
}

func CancelBooking(storage config.Database, bookingID uuid.UUID) (config.Booking, error) {
	result, err := storage.CancelBookingByBookingID(bookingID)
	return result, err
}
