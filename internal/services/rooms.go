package services

import (
	"avito_tech_backend/internal/config"
	"time"

	"github.com/google/uuid"
)

func RoomExist(storage config.Database, roomID uuid.UUID) bool {
	_, err := storage.GetRoom(roomID)
	return err == nil
}

func RoomsList(storage config.Database) ([]config.Room, error) {
	result, err := storage.GetRoomsList()
	return result, err
}

func RoomsCreate(storage config.Database, user config.User, room config.Room) (config.Room, error) {
	room.ID, _ = uuid.NewUUID()
	room.CreatedAt = time.Now().Format(time.RFC3339)

	err := storage.CreateRoom(user.ID, room)

	if err != nil {
		return config.Room{}, err
	}
	return room, nil
}
