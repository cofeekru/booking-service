package services_test

import (
	"errors"
	"testing"

	"avito_tech_backend/internal/config"
	"avito_tech_backend/internal/services"

	"avito_tech_backend/internal/mocks"

	"github.com/google/uuid"
)

func TestRoomExist(t *testing.T) {
	tests := []struct {
		name       string
		setup      func() *mocks.MockDatabase
		roomID     uuid.UUID
		wantExists bool
	}{
		{
			name: "Room exists - no error from storage",
			setup: func() *mocks.MockDatabase {
				mock := &mocks.MockDatabase{
					GetRoomFunc: func(roomID uuid.UUID) (config.Room, error) {
						return config.Room{ID: roomID}, nil
					},
				}
				return mock
			},
			roomID:     uuid.New(),
			wantExists: true,
		},
		{
			name: "Room does not exist - error from storage",
			setup: func() *mocks.MockDatabase {
				mock := &mocks.MockDatabase{
					GetRoomFunc: func(roomID uuid.UUID) (config.Room, error) {
						return config.Room{}, errors.New("room not found")
					},
				}
				return mock
			},
			roomID:     uuid.New(),
			wantExists: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := tt.setup()

			exists := services.RoomExist(mockDB, tt.roomID)

			if exists != tt.wantExists {
				t.Errorf("RoomExist() = %v, want %v", exists, tt.wantExists)
			}
		})
	}
}

func TestRoomsList(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *mocks.MockDatabase
		wantNil bool
		wantErr bool
	}{
		{
			name: "Successful rooms list retrieval",
			setup: func() *mocks.MockDatabase {
				mock := &mocks.MockDatabase{
					GetRoomsListFunc: func() ([]config.Room, error) {
						return []config.Room{{ID: uuid.New()}}, nil
					},
				}
				return mock
			},
			wantNil: false,
			wantErr: false,
		},
		{
			name: "Error retrieving rooms list",
			setup: func() *mocks.MockDatabase {
				mock := &mocks.MockDatabase{
					GetRoomsListFunc: func() ([]config.Room, error) {
						return nil, errors.New("database error")
					},
				}
				return mock
			},
			wantNil: true,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := tt.setup()

			result, err := services.RoomsList(mockDB)

			if (err != nil) != tt.wantErr {
				t.Errorf("RoomsList() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && result == nil {
				t.Error("RoomsList() returned nil result, want non-nil")
			}
		})
	}
}

func TestRoomsCreate(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *mocks.MockDatabase
		user    config.User
		room    config.Room
		wantErr bool
	}{
		{
			name: "Successful room creation",
			setup: func() *mocks.MockDatabase {
				mock := &mocks.MockDatabase{
					CreateRoomFunc: func(userID uuid.UUID, room config.Room) error {
						return nil
					},
				}
				return mock
			},
			user:    config.User{ID: uuid.New()},
			room:    config.Room{Name: "Test Room"},
			wantErr: false,
		},
		{
			name: "Error during room creation",
			setup: func() *mocks.MockDatabase {
				mock := &mocks.MockDatabase{
					CreateRoomFunc: func(userID uuid.UUID, room config.Room) error {
						return errors.New("database error")
					},
				}
				return mock
			},
			user:    config.User{ID: uuid.New()},
			room:    config.Room{Name: "Test Room"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := tt.setup()

			_, err := services.RoomsCreate(mockDB, tt.user, tt.room)

			if (err != nil) != tt.wantErr {
				t.Errorf("RoomsCreate() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}
