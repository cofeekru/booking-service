package database

import (
	"avito_tech_backend/internal/config"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

func (storage *Storage) GetSlotsList(roomID uuid.UUID, dateRequest time.Time) ([]config.Slot, error) {
	rows, err := storage.DB.Query(`
		SELECT slot_id, room_id, start_slot, end_slot 
		FROM slots 
		WHERE (booking_status = 'cancelled' OR booking_status IS NULL) AND EXTRACT(DAY from start_slot) = $1 AND EXTRACT(MONTH from start_slot) = $2 AND EXTRACT(YEAR from start_slot) = $3`,
		dateRequest.Day(), dateRequest.Month(), dateRequest.Year())

	if err != nil {
		return []config.Slot{}, err
	}

	var slots []config.Slot
	for rows.Next() {
		var slot config.Slot

		if err := rows.Scan(&slot.ID, &slot.RoomID, &slot.Start, &slot.End); err != nil {
			return nil, err
		}
		slots = append(slots, slot)
	}
	return slots, nil
}

func (storage *Storage) GetSlotBySlotID(slotID uuid.UUID) (config.Slot, error) {
	var slot config.Slot
	err := storage.DB.QueryRow(`
		SELECT slot_id, room_id, start_slot, end_slot 
		FROM slots WHERE slot_id = $1`, slotID).
		Scan(&slot.ID, &slot.RoomID, &slot.Start, &slot.End)

	if err != nil {
		slog.Error(err.Error())
		return config.Slot{}, err
	}
	return slot, nil
}

func (storage *Storage) CreateSlot(slot config.Slot, startSlot time.Time, endSlot time.Time) error {
	_, err := storage.DB.Exec(`
		INSERT INTO slots (slot_id, room_id, start_slot, end_slot) 
		VALUES ($1, $2, $3, $4)`,
		slot.ID, slot.RoomID, startSlot, endSlot)
	return err
}
