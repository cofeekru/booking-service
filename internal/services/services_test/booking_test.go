package services_test

import (
	"errors"
	"testing"

	"avito_tech_backend/internal/config"
	"avito_tech_backend/internal/mocks"
	"avito_tech_backend/internal/services"

	"github.com/google/uuid"
)

func TestCreateBooking(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *mocks.MockDatabase
		wantErr bool
	}{
		{
			name: "Successful booking creation",
			setup: func() *mocks.MockDatabase {
				mock := &mocks.MockDatabase{
					CreateBookingFunc: func(booking *config.Booking) error {
						return nil
					},
				}
				return mock
			},
			wantErr: false,
		},
		{
			name: "Error during booking creation",
			setup: func() *mocks.MockDatabase {
				mock := &mocks.MockDatabase{
					CreateBookingFunc: func(booking *config.Booking) error {
						return errors.New("database error")
					},
				}
				return mock
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := tt.setup()
			booking := &config.Booking{}

			err := services.CreateBooking(mockDB, booking)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateBooking() error = %v, wantErr %v", err, tt.wantErr)
			}

			if booking.Status != "active" {
				t.Errorf("CreateBooking() status = %s, want 'active'", booking.Status)
			}
		})
	}
}

func TestBookingsList(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *mocks.MockDatabase
		wantLen int
		wantErr bool
	}{
		{
			name: "Successful admin bookings list retrieval",
			setup: func() *mocks.MockDatabase {
				mock := &mocks.MockDatabase{
					GetAdminBookingListFunc: func(pagination *config.Pagination, user config.User) ([]config.Booking, error) {
						return []config.Booking{{ID: uuid.New()}}, nil
					},
				}
				return mock
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "Error retrieving admin bookings",
			setup: func() *mocks.MockDatabase {
				mock := &mocks.MockDatabase{
					GetAdminBookingListFunc: func(pagination *config.Pagination, user config.User) ([]config.Booking, error) {
						return nil, errors.New("database error")
					},
				}
				return mock
			},
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := tt.setup()
			pagination := &config.Pagination{Page: 1, PageSize: 5}
			user := config.User{ID: uuid.New()}

			result, err := services.BookingsList(mockDB, pagination, user)

			if (err != nil) != tt.wantErr {
				t.Errorf("BookingsList() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(result) != tt.wantLen {
				t.Errorf("BookingsList() returned %d bookings, want %d", len(result), tt.wantLen)
			}
		})
	}
}

func TestBookingsMy(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *mocks.MockDatabase
		wantLen int
		wantErr bool
	}{
		{
			name: "Successful user bookings retrieval",
			setup: func() *mocks.MockDatabase {
				mock := &mocks.MockDatabase{
					GetUserBookingListFunc: func(user config.User) ([]config.Booking, error) {
						return []config.Booking{{ID: uuid.New()}}, nil
					},
				}
				return mock
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "Error retrieving user bookings",
			setup: func() *mocks.MockDatabase {
				mock := &mocks.MockDatabase{
					GetUserBookingListFunc: func(user config.User) ([]config.Booking, error) {
						return nil, errors.New("database error")
					},
				}
				return mock
			},
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := tt.setup()
			user := config.User{ID: uuid.New()}

			result, err := services.BookingsMy(mockDB, user)

			if (err != nil) != tt.wantErr {
				t.Errorf("BookingsMy() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(result) != tt.wantLen {
				t.Errorf("BookingsMy() returned %d bookings, want %d", len(result), tt.wantLen)
			}
		})
	}
}

func TestBookingExist(t *testing.T) {
	tests := []struct {
		name       string
		setup      func() *mocks.MockDatabase
		bookingID  uuid.UUID
		wantExists bool
	}{
		{
			name: "Booking exists — successful retrieval (no error)",
			setup: func() *mocks.MockDatabase {
				mock := &mocks.MockDatabase{
					GetBookingByBookingIDFunc: func(bookingID uuid.UUID) (config.Booking, error) {
						return config.Booking{ID: bookingID}, nil
					},
				}
				return mock
			},
			bookingID:  uuid.New(),
			wantExists: true,
		},
		{
			name: "Booking does not exist — returns error (e.g., not found)",
			setup: func() *mocks.MockDatabase {
				mock := &mocks.MockDatabase{
					GetBookingByBookingIDFunc: func(bookingID uuid.UUID) (config.Booking, error) {
						return config.Booking{}, errors.New("booking not found")
					},
				}
				return mock
			},
			bookingID:  uuid.New(),
			wantExists: false,
		},
		{
			name: "Database error during booking lookup — treats as non‑existent",
			setup: func() *mocks.MockDatabase {
				mock := &mocks.MockDatabase{
					GetBookingByBookingIDFunc: func(bookingID uuid.UUID) (config.Booking, error) {
						return config.Booking{}, errors.New("database connection failed")
					},
				}
				return mock
			},
			bookingID:  uuid.New(),
			wantExists: false,
		},
		{
			name: "Empty booking ID — still checks existence (edge case)",
			setup: func() *mocks.MockDatabase {
				mock := &mocks.MockDatabase{
					GetBookingByBookingIDFunc: func(bookingID uuid.UUID) (config.Booking, error) {
						if bookingID == uuid.Nil {
							return config.Booking{}, errors.New("invalid booking ID")
						}
						return config.Booking{ID: bookingID}, nil
					},
				}
				return mock
			},
			bookingID:  uuid.Nil,
			wantExists: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := tt.setup()

			exists := services.BookingExist(mockDB, tt.bookingID)

			if exists != tt.wantExists {
				t.Errorf(
					"BookingExist() = %v, want %v (test: %q)",
					exists,
					tt.wantExists,
					tt.name,
				)
			}
		})
	}
}

func TestBookingValid(t *testing.T) {
	tests := []struct {
		name       string
		setup      func() *mocks.MockDatabase
		firstUUID  uuid.UUID
		secondUUID uuid.UUID
		bookingID  uuid.UUID
		userID     uuid.UUID
		wantValid  bool
	}{
		{
			name: "Booking is valid — user IDs match",
			setup: func() *mocks.MockDatabase {
				mock := &mocks.MockDatabase{
					GetBookingByBookingIDFunc: func(bookingID uuid.UUID) (config.Booking, error) {
						return config.Booking{
							ID:     bookingID,
							UserID: uuid.Nil,
						}, nil
					},
				}
				return mock
			},
			bookingID: uuid.New(),
			userID:    uuid.Nil,
			wantValid: true,
		},
		{
			name: "Booking is not valid — user IDs don't match",
			setup: func() *mocks.MockDatabase {
				mock := &mocks.MockDatabase{
					GetBookingByBookingIDFunc: func(bookingID uuid.UUID) (config.Booking, error) {
						return config.Booking{
							ID:     bookingID,
							UserID: uuid.New(),
						}, nil
					},
				}
				return mock
			},
			bookingID: uuid.New(),
			userID:    uuid.New(),
			wantValid: false,
		},
		{
			name: "Booking not found — returns false",
			setup: func() *mocks.MockDatabase {
				mock := &mocks.MockDatabase{
					GetBookingByBookingIDFunc: func(bookingID uuid.UUID) (config.Booking, error) {
						return config.Booking{}, errors.New("booking not found")
					},
				}
				return mock
			},
			bookingID: uuid.New(),
			userID:    uuid.New(),
			wantValid: false,
		},
		{
			name: "Database error during booking lookup — returns false",
			setup: func() *mocks.MockDatabase {
				mock := &mocks.MockDatabase{
					GetBookingByBookingIDFunc: func(bookingID uuid.UUID) (config.Booking, error) {
						return config.Booking{}, errors.New("database error")
					},
				}
				return mock
			},
			bookingID: uuid.New(),
			userID:    uuid.New(),
			wantValid: false,
		},
		{
			name: "Zero UUID booking ID — returns false",
			setup: func() *mocks.MockDatabase {
				mock := &mocks.MockDatabase{
					GetBookingByBookingIDFunc: func(bookingID uuid.UUID) (config.Booking, error) {
						if bookingID == uuid.Nil {
							return config.Booking{}, errors.New("invalid booking ID")
						}
						return config.Booking{ID: bookingID, UserID: uuid.New()}, nil
					},
				}
				return mock
			},
			bookingID: uuid.Nil,
			userID:    uuid.New(),
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := tt.setup()

			valid := services.BookingValid(mockDB, tt.bookingID, tt.userID)

			if valid != tt.wantValid {
				t.Errorf(
					"BookingValid() = %v, want %v (test: %q)",
					valid,
					tt.wantValid,
					tt.name,
				)
			}
		})
	}
}

func TestCancelBooking(t *testing.T) {
	tests := []struct {
		name      string
		setup     func() *mocks.MockDatabase
		bookingID uuid.UUID
		wantErr   bool
	}{
		{
			name: "Successful booking cancellation",
			setup: func() *mocks.MockDatabase {
				mock := &mocks.MockDatabase{
					CancelBookingByBookingIDFunc: func(bookingID uuid.UUID) (config.Booking, error) {
						return config.Booking{ID: bookingID, Status: "cancelled"}, nil
					},
				}
				return mock
			},
			bookingID: uuid.New(),
			wantErr:   false,
		},
		{
			name: "Error during booking cancellation",
			setup: func() *mocks.MockDatabase {
				mock := &mocks.MockDatabase{
					CancelBookingByBookingIDFunc: func(bookingID uuid.UUID) (config.Booking, error) {
						return config.Booking{}, errors.New("cancellation failed")
					},
				}
				return mock
			},
			bookingID: uuid.New(),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := tt.setup()

			result, err := services.CancelBooking(mockDB, tt.bookingID)

			if (err != nil) != tt.wantErr {
				t.Errorf("CancelBooking() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && result.ID != tt.bookingID {
				t.Errorf("CancelBooking() returned booking with ID %v, want %v", result.ID, tt.bookingID)
			}
		})
	}
}
