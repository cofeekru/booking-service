package services

import (
	"avito_tech_backend/internal/config"

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

func RoomsCreate(storage config.Database, user config.User, room *config.Room) error {
	room.ID, _ = uuid.NewUUID()

	err := storage.CreateRoom(user.ID, room)
	return err
}
