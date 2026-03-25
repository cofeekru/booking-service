package mocks

import (
	"avito_tech_backend/internal/config"
	"time"

	"github.com/google/uuid"
)

type MockDatabase struct {
	CreateScheduleFunc      func(schedule *config.Schedule, roomID uuid.UUID, user config.User) error
	CreateSlotFunc          func(slot config.Slot, startSlot time.Time, endSlot time.Time) error
	GetSlotsListFunc        func(roomID uuid.UUID, dateRequest time.Time) ([]config.Slot, error)
	GetSlotBySlotIDFunc     func(slotID uuid.UUID) (config.Slot, error)
	GetScheduleByRoomIDFunc func(roomID uuid.UUID) (config.Schedule, error)

	CreateBookingFunc            func(booking *config.Booking) error
	GetAdminBookingListFunc      func(pagination *config.Pagination, user config.User) ([]config.Booking, error)
	GetUserBookingListFunc       func(user config.User) ([]config.Booking, error)
	GetBookingByBookingIDFunc    func(bookingID uuid.UUID) (config.Booking, error)
	GetBookingBySlotIDFunc       func(slotID uuid.UUID) (config.Booking, error)
	CancelBookingByBookingIDFunc func(bookingID uuid.UUID) (config.Booking, error)

	CreateRoomFunc   func(userID uuid.UUID, room *config.Room) error
	GetRoomFunc      func(roomID uuid.UUID) (config.Room, error)
	GetRoomsListFunc func() ([]config.Room, error)
}

func (m *MockDatabase) CreateSchedule(schedule *config.Schedule, roomID uuid.UUID, user config.User) error {
	if m.CreateScheduleFunc != nil {
		return m.CreateScheduleFunc(schedule, roomID, user)
	}
	return nil
}

func (m *MockDatabase) CreateSlot(slot config.Slot, startSlot time.Time, endSlot time.Time) error {
	if m.CreateSlotFunc != nil {
		return m.CreateSlotFunc(slot, startSlot, endSlot)
	}
	return nil
}

func (m *MockDatabase) GetSlotsList(roomID uuid.UUID, dateRequest time.Time) ([]config.Slot, error) {
	if m.GetSlotsListFunc != nil {
		return m.GetSlotsListFunc(roomID, dateRequest)
	}
	return nil, nil
}

func (m *MockDatabase) GetSlotBySlotID(slotID uuid.UUID) (config.Slot, error) {
	if m.GetSlotBySlotIDFunc != nil {
		return m.GetSlotBySlotIDFunc(slotID)
	}
	return config.Slot{}, nil
}

func (m *MockDatabase) GetScheduleByRoomID(roomID uuid.UUID) (config.Schedule, error) {
	if m.GetScheduleByRoomIDFunc != nil {
		return m.GetScheduleByRoomIDFunc(roomID)
	}
	return config.Schedule{}, nil
}

func (m *MockDatabase) CreateBooking(booking *config.Booking) error {
	if m.CreateBookingFunc != nil {
		return m.CreateBookingFunc(booking)
	}
	return nil
}

func (m *MockDatabase) GetAdminBookingList(pagination *config.Pagination, user config.User) ([]config.Booking, error) {
	if m.GetAdminBookingListFunc != nil {
		return m.GetAdminBookingListFunc(pagination, user)
	}
	return nil, nil
}

func (m *MockDatabase) GetUserBookingList(user config.User) ([]config.Booking, error) {
	if m.GetUserBookingListFunc != nil {
		return m.GetUserBookingListFunc(user)
	}
	return nil, nil
}

func (m *MockDatabase) GetBookingByBookingID(bookingID uuid.UUID) (config.Booking, error) {
	if m.GetBookingByBookingIDFunc != nil {
		return m.GetBookingByBookingIDFunc(bookingID)
	}
	return config.Booking{}, nil
}

func (m *MockDatabase) GetBookingBySlotID(slotID uuid.UUID) (config.Booking, error) {
	if m.GetBookingBySlotIDFunc != nil {
		return m.GetBookingBySlotIDFunc(slotID)
	}
	return config.Booking{}, nil
}

func (m *MockDatabase) CancelBookingByBookingID(bookingID uuid.UUID) (config.Booking, error) {
	if m.CancelBookingByBookingIDFunc != nil {
		return m.CancelBookingByBookingIDFunc(bookingID)
	}
	return config.Booking{}, nil
}

func (m *MockDatabase) CreateRoom(userID uuid.UUID, room *config.Room) error {
	if m.CreateRoomFunc != nil {
		return m.CreateRoomFunc(userID, room)
	}
	return nil
}

func (m *MockDatabase) GetRoom(roomID uuid.UUID) (config.Room, error) {
	if m.GetRoomFunc != nil {
		return m.GetRoomFunc(roomID)
	}
	return config.Room{}, nil
}

func (m *MockDatabase) GetRoomsList() ([]config.Room, error) {
	if m.GetRoomsListFunc != nil {
		return m.GetRoomsListFunc()
	}
	return nil, nil
}
