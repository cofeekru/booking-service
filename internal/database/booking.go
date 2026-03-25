package database

import (
	"avito_tech_backend/internal/config"

	"time"

	"github.com/google/uuid"
)

func (storage *Storage) GetBookingBySlotID(slotID uuid.UUID) (config.Booking, error) {
	var booking config.Booking
	err := storage.DB.QueryRow(`
		SELECT booking_id, booking_status, booking_user_id 
		FROM slots 
		WHERE slot_id = $1`, slotID).
		Scan(&booking.ID, &booking.Status, &booking.UserID)

	if err != nil {
		return config.Booking{}, err
	}
	return booking, nil
}

func (storage *Storage) CreateBooking(booking *config.Booking) error {
	_, err := storage.DB.Exec(`
		UPDATE slots 
		SET booking_id = $1, booking_status = $2, booking_user_id = $3 
		WHERE slot_id = $4`, booking.ID, booking.Status, booking.UserID, booking.SlotID)

	return err
}

func (storage *Storage) GetAdminBookingList(pagination *config.Pagination, user config.User) ([]config.Booking, error) {
	rows, err := storage.DB.Query(`
		SELECT booking_id, slot_id, booking_user_id, booking_status 
		FROM slots INNER JOIN rooms ON
			slots.room_id = rooms.room_id
		WHERE booking_status = 'active' AND user_id = $1
		LIMIT $2 OFFSET $3`, user.ID, pagination.PageSize, (pagination.Page-1)*pagination.PageSize)

	if err != nil {
		return nil, err
	}

	var bookings []config.Booking
	for rows.Next() {
		var booking config.Booking

		if err := rows.Scan(&booking.ID, &booking.SlotID, &booking.UserID, &booking.Status); err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}

	pagination.Total = len(bookings)
	return bookings, nil
}

func (storage *Storage) GetUserBookingList(user config.User) ([]config.Booking, error) {
	rows, err := storage.DB.Query(`
		SELECT booking_id, slot_id, booking_user_id, booking_status 
		FROM slots 
		WHERE booking_status = 'active' AND booking_user_id = $1 AND start_slot >= $2`, user.ID, time.Now().Format(time.RFC3339))

	if err != nil {
		return nil, err
	}

	var bookings []config.Booking
	for rows.Next() {
		var booking config.Booking

		if err := rows.Scan(&booking.ID, &booking.SlotID, &booking.UserID, &booking.Status); err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}

	return bookings, nil
}

func (storage *Storage) GetBookingByBookingID(bookingID uuid.UUID) (config.Booking, error) {
	var booking config.Booking
	err := storage.DB.QueryRow(`
		SELECT booking_id, slot_id, booking_status, booking_user_id 
		FROM slots 
		WHERE booking_id = $1`, bookingID).
		Scan(&booking.ID, &booking.SlotID, &booking.Status, &booking.UserID)

	if err != nil {
		return config.Booking{}, err
	}
	return booking, nil
}
func (storage *Storage) CancelBookingByBookingID(bookingID uuid.UUID) (config.Booking, error) {
	_, err := storage.DB.Exec(`
		UPDATE slots 
		SET booking_status = $1 
		WHERE booking_id = $2`, "cancelled", bookingID)
	if err != nil {
		return config.Booking{}, err
	}

	var booking config.Booking
	booking, err = storage.GetBookingByBookingID(bookingID)
	if err != nil {
		return config.Booking{}, err
	}
	return booking, nil
}
