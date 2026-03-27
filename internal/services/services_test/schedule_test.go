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

// Тесты для ScheduleExist
func TestScheduleExist(t *testing.T) {
	tests := []struct {
		name       string
		setup      func() *mocks.MockDatabase
		schedule   *config.Schedule
		wantExists bool
	}{
		{
			name: "Schedule exists — successful retrieval",
			setup: func() *mocks.MockDatabase {
				mock := &mocks.MockDatabase{
					GetScheduleByRoomIDFunc: func(roomID uuid.UUID) (config.Schedule, error) {
						return config.Schedule{RoomID: roomID}, nil
					},
				}
				return mock
			},
			schedule: &config.Schedule{
				RoomID: uuid.New(),
			},
			wantExists: true,
		},
		{
			name: "Schedule does not exist — returns error",
			setup: func() *mocks.MockDatabase {
				mock := &mocks.MockDatabase{
					GetScheduleByRoomIDFunc: func(roomID uuid.UUID) (config.Schedule, error) {
						return config.Schedule{}, errors.New("schedule not found")
					},
				}
				return mock
			},
			schedule: &config.Schedule{
				RoomID: uuid.New(),
			},
			wantExists: false,
		},
		{
			name: "Database error during schedule lookup",
			setup: func() *mocks.MockDatabase {
				mock := &mocks.MockDatabase{
					GetScheduleByRoomIDFunc: func(roomID uuid.UUID) (config.Schedule, error) {
						return config.Schedule{}, errors.New("database error")
					},
				}
				return mock
			},
			schedule: &config.Schedule{
				RoomID: uuid.New(),
			},
			wantExists: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := tt.setup()

			exists := services.ScheduleExist(mockDB, tt.schedule)

			if exists != tt.wantExists {
				t.Errorf("ScheduleExist() = %v, want %v (test: %q)", exists, tt.wantExists, tt.name)
			}
		})
	}
}

func TestValidSchedule(t *testing.T) {
	tests := []struct {
		name      string
		schedule  config.Schedule
		wantValid bool
	}{
		{
			name: "Valid schedule — correct days of week (1–7)",
			schedule: config.Schedule{
				DaysOfWeek: []int{1, 3, 5},
			},
			wantValid: true,
		},
		{
			name: "Invalid — empty days of week",
			schedule: config.Schedule{
				DaysOfWeek: []int{},
			},
			wantValid: false,
		},
		{
			name: "Invalid — more than 7 days",
			schedule: config.Schedule{
				DaysOfWeek: []int{1, 2, 3, 4, 5, 6, 7, 8},
			},
			wantValid: false,
		},
		{
			name: "Invalid — days out of range (0)",
			schedule: config.Schedule{
				DaysOfWeek: []int{0, 1, 2},
			},
			wantValid: false,
		},
		{
			name: "Invalid — duplicate days",
			schedule: config.Schedule{
				DaysOfWeek: []int{1, 1, 2},
			},
			wantValid: false,
		},
		{
			name: "Invalid — days out of range (8)",
			schedule: config.Schedule{
				DaysOfWeek: []int{8, 1},
			},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := services.ValidSchedule(tt.schedule)

			if valid != tt.wantValid {
				t.Errorf("ValidSchedule() = %v, want %v (test: %q)", valid, tt.wantValid, tt.name)
			}
		})
	}
}

func TestScheduleCreate(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *mocks.MockDatabase
		schedule *config.Schedule
		roomID   uuid.UUID
		user     config.User
		wantErr  bool
	}{
		{
			name: "Successful schedule creation",
			setup: func() *mocks.MockDatabase {
				mock := &mocks.MockDatabase{
					CreateScheduleFunc: func(schedule *config.Schedule, roomID uuid.UUID, user config.User) error {
						return nil
					},
					CreateSlotFunc: func(slot config.Slot, startSlot time.Time, endSlot time.Time) error {
						return nil
					},
				}
				return mock
			},
			schedule: &config.Schedule{
				StartTime:  "09:00",
				EndTime:    "10:00",
				DaysOfWeek: []int{1, 3},
				RoomID:     uuid.New(),
			},
			roomID:  uuid.New(),
			user:    config.User{ID: uuid.New()},
			wantErr: false,
		},
		{
			name: "Error creating schedule in storage",
			setup: func() *mocks.MockDatabase {
				mock := &mocks.MockDatabase{
					CreateScheduleFunc: func(schedule *config.Schedule, roomID uuid.UUID, user config.User) error {
						return errors.New("storage error")
					},
				}
				return mock
			},
			schedule: &config.Schedule{
				StartTime:  "09:00",
				EndTime:    "10:00",
				DaysOfWeek: []int{1},
				RoomID:     uuid.New(),
			},
			roomID:  uuid.New(),
			user:    config.User{ID: uuid.New()},
			wantErr: true,
		},
		{
			name: "Error creating slot in storage",
			setup: func() *mocks.MockDatabase {
				mock := &mocks.MockDatabase{
					CreateScheduleFunc: func(schedule *config.Schedule, roomID uuid.UUID, user config.User) error {
						return nil
					},
					CreateSlotFunc: func(slot config.Slot, startSlot time.Time, endSlot time.Time) error {
						return errors.New("slot creation error")
					},
				}
				return mock
			},
			schedule: &config.Schedule{
				StartTime:  "09:00",
				EndTime:    "11:00",
				DaysOfWeek: []int{1},
				RoomID:     uuid.New(),
			},
			roomID:  uuid.New(),
			user:    config.User{ID: uuid.New()},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := tt.setup()

			err := services.ScheduleCreate(mockDB, tt.schedule, tt.roomID, tt.user)

			if (err != nil) != tt.wantErr {
				t.Errorf("ScheduleCreate() error = %v, wantErr %v (test: %q)", err, tt.wantErr, tt.name)
			}

			if err == nil && tt.schedule.ID.String() == "" {
				t.Error("ScheduleCreate() did not generate UUID for schedule ID")
			}
		})
	}
}
