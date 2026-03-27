package database

import (
	"avito_tech_backend/internal/config"

	"github.com/google/uuid"
)

func (storage *Storage) CreateRoom(userID uuid.UUID, room config.Room) error {
	_, err := storage.DB.Exec(`
		INSERT INTO rooms (room_id, user_id, name, description, capacity, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6)`,
		room.ID, userID, room.Name, room.Description, room.Capacity, room.CreatedAt)
	return err
}

func (storage *Storage) GetRoom(roomID uuid.UUID) (config.Room, error) {
	var room config.Room
	err := storage.DB.QueryRow(`
		SELECT room_id, name, description, capacity, created_at 
		FROM rooms 
		WHERE room_id = $1`, roomID.String()).
		Scan(&room.ID, &room.Name, &room.Description, &room.Capacity, &room.CreatedAt)

	return room, err
}

func (storage *Storage) GetRoomsList() ([]config.Room, error) {
	rows, err := storage.DB.Query(`
		SELECT room_id, name, description, capacity, created_at 
		FROM rooms`)

	if err != nil {
		return nil, err
	}

	var rooms []config.Room
	for rows.Next() {
		var room config.Room

		if err := rows.Scan(&room.ID, &room.Name, &room.Description, &room.Capacity, &room.CreatedAt); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}
	return rooms, nil
}
