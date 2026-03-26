package services

import (
	"avito_tech_backend/internal/config"
	"log/slog"

	"time"

	"github.com/google/uuid"
)

func SlotsList(storage config.Database, roomID uuid.UUID, dateRequestString string) ([]config.Slot, error) {
	date, err := time.Parse("2006-01-02", dateRequestString)
	if err != nil {
		slog.Error(err.Error())
		return []config.Slot{}, err
	}

	result, err := storage.GetSlotsList(roomID, date)
	if err != nil {
		slog.Error(err.Error())
		return []config.Slot{}, err
	}
	return result, err
}

func SlotExist(storage config.Database, slotID uuid.UUID) bool {
	_, err := storage.GetSlotBySlotID(slotID)
	return err == nil
}
func SlotBooked(storage config.Database, slotID uuid.UUID) bool {
	result, _ := storage.GetBookingBySlotID(slotID)

	return result.Status == "active"
}

func SlotInPast(storage config.Database, slotID uuid.UUID) bool {
	result, _ := storage.GetSlotBySlotID(slotID)
	timeSlot, _ := time.Parse(time.RFC3339, result.Start)
	return timeSlot.Before(time.Now())
}
