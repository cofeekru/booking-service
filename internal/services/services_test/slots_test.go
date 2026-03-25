package services_test

import (
	"errors"
	"testing"
	"time"

	"avito_tech_backend/internal/config"
	"avito_tech_backend/internal/mocks"
	"avito_tech_backend/internal/services"

	"github.com/google/uuid"
)

// Тесты для SlotsList
func TestSlotsList(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() *mocks.MockDatabase
		roomID      uuid.UUID
		dateRequest string
		wantLen     int
		wantErr     bool
	}{
		{
			name: "Successful slots list retrieval with valid date",
			setup: func() *mocks.MockDatabase {
				mock := &mocks.MockDatabase{
					GetSlotsListFunc: func(roomID uuid.UUID, date time.Time) ([]config.Slot, error) {
						return []config.Slot{
							{ID: uuid.New(), RoomID: roomID, Start: date.Format(time.RFC3339)},
						}, nil
					},
				}
				return mock
			},
			roomID:      uuid.New(),
			dateRequest: "2023-12-25",
			wantLen:     1,
			wantErr:     false,
		},
		{
			name: "Invalid date format — returns error",
			setup: func() *mocks.MockDatabase {
				return &mocks.MockDatabase{} // не используем БД, так как ошибка на этапе парсинга даты
			},
			roomID:      uuid.New(),
			dateRequest: "invalid-date-format",
			wantLen:     0,
			wantErr:     true,
		},
		{
			name: "Error retrieving slots from storage",
			setup: func() *mocks.MockDatabase {
				mock := &mocks.MockDatabase{
					GetSlotsListFunc: func(roomID uuid.UUID, date time.Time) ([]config.Slot, error) {
						return nil, errors.New("database error")
					},
				}
				return mock
			},
			roomID:      uuid.New(),
			dateRequest: "2023-12-25",
			wantLen:     0,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := tt.setup()

			result, err := services.SlotsList(mockDB, tt.roomID, tt.dateRequest)

			if (err != nil) != tt.wantErr {
				t.Errorf("SlotsList() error = %v, wantErr %v (test: %q)", err, tt.wantErr, tt.name)
			}

			if len(result) != tt.wantLen {
				t.Errorf("SlotsList() returned %d slots, want %d (test: %q)", len(result), tt.wantLen, tt.name)
			}
		})
	}
}

// Тесты для SlotExist
func TestSlotExist(t *testing.T) {
	tests := []struct {
		name      string
		setup     func() *mocks.MockDatabase
		slotID    uuid.UUID
		wantExist bool
	}{
		{
			name: "Slot exists — successful retrieval",
			setup: func() *mocks.MockDatabase {
				mock := &mocks.MockDatabase{
					GetSlotBySlotIDFunc: func(slotID uuid.UUID) (config.Slot, error) {
						return config.Slot{ID: slotID}, nil
					},
				}
				return mock
			},
			slotID:    uuid.New(),
			wantExist: true,
		},
		{
			name: "Slot does not exist — returns error",
			setup: func() *mocks.MockDatabase {
				mock := &mocks.MockDatabase{
					GetSlotBySlotIDFunc: func(slotID uuid.UUID) (config.Slot, error) {
						return config.Slot{}, errors.New("slot not found")
					},
				}
				return mock
			},
			slotID:    uuid.New(),
			wantExist: false,
		},
		{
			name: "Database error during slot lookup",
			setup: func() *mocks.MockDatabase {
				mock := &mocks.MockDatabase{
					GetSlotBySlotIDFunc: func(slotID uuid.UUID) (config.Slot, error) {
						return config.Slot{}, errors.New("database error")
					},
				}
				return mock
			},
			slotID:    uuid.New(),
			wantExist: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := tt.setup()

			exists := services.SlotExist(mockDB, tt.slotID)

			if exists != tt.wantExist {
				t.Errorf("SlotExist() = %v, want %v (test: %q)", exists, tt.wantExist, tt.name)
			}
		})
	}
}

// Тесты для SlotBooked
func TestSlotBooked(t *testing.T) {
	tests := []struct {
		name       string
		setup      func() *mocks.MockDatabase
		slotID     uuid.UUID
		wantBooked bool
	}{
		{
			name: "Slot is booked — status is active",
			setup: func() *mocks.MockDatabase {
				mock := &mocks.MockDatabase{
					GetBookingBySlotIDFunc: func(slotID uuid.UUID) (config.Booking, error) {
						return config.Booking{Status: "active"}, nil
					},
				}
				return mock
			},
			slotID:     uuid.New(),
			wantBooked: true,
		},
		{
			name: "Slot is not booked — status is not active",
			setup: func() *mocks.MockDatabase {
				mock := &mocks.MockDatabase{
					GetBookingBySlotIDFunc: func(slotID uuid.UUID) (config.Booking, error) {
						return config.Booking{Status: "cancelled"}, nil
					},
				}
				return mock
			},
			slotID:     uuid.New(),
			wantBooked: false,
		},
		{
			name: "No booking found for slot",
			setup: func() *mocks.MockDatabase {
				mock := &mocks.MockDatabase{
					GetBookingBySlotIDFunc: func(slotID uuid.UUID) (config.Booking, error) {
						return config.Booking{}, nil // пустое бронирование
					},
				}
				return mock
			},
			slotID:     uuid.New(),
			wantBooked: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := tt.setup()

			booked := services.SlotBooked(mockDB, tt.slotID)

			if booked != tt.wantBooked {
				t.Errorf("SlotBooked() = %v, want %v (test: %q)", booked, tt.wantBooked, tt.name)
			}
		})
	}
}

// Тесты для SlotInPast
func TestSlotInPast(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *mocks.MockDatabase
		slotID   uuid.UUID
		wantPast bool
	}{
		{
			name: "Slot is in the past",
			setup: func() *mocks.MockDatabase {
				pastTime := time.Now().Add(-time.Hour)
				mock := &mocks.MockDatabase{
					GetSlotBySlotIDFunc: func(slotID uuid.UUID) (config.Slot, error) {
						return config.Slot{Start: pastTime.Format(time.RFC3339)}, nil
					},
				}
				return mock
			},
			slotID:   uuid.New(),
			wantPast: true,
		},
		{
			name: "Slot is in the future",
			setup: func() *mocks.MockDatabase {
				futureTime := time.Now().Add(time.Hour)
				mock := &mocks.MockDatabase{
					GetSlotBySlotIDFunc: func(slotID uuid.UUID) (config.Slot, error) {
						return config.Slot{Start: futureTime.Format(time.RFC3339)}, nil
					},
				}
				return mock
			},
			slotID:   uuid.New(),
			wantPast: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := tt.setup()

			past := services.SlotInPast(mockDB, tt.slotID)

			if past != tt.wantPast {
				t.Errorf("SlotInPast() = %v, want %v (test: %q)", past, tt.wantPast, tt.name)
			}
		})
	}
}
