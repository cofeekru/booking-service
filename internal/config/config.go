package config

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

func MustLoad() *Config {
	var cfg Config
	godotenv.Load(".env")
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("Cannot read config: %v", err)
	}
	SecretKey = []byte(cfg.SECRET_KEY)

	return &cfg
}

type Config struct {
	HTTPServer
	DB_HOST     string `env:"DB_HOST"`
	DB_NAME     string `env:"DB_NAME"`
	DB_PORT     string `env:"DB_PORT"`
	DB_USER     string `env:"DB_USER"`
	DB_PASSWORD string `env:"DB_PASSWORD"`
	DB_SSL_MODE string `env:"DB_SSL_MODE"`
	SECRET_KEY  string `env:"SECRET_KEY"`
}

type HTTPServer struct {
	Host        string        `env:"HOST"`
	Timeout     time.Duration `env:"TIMEOUT"`
	IdleTimeout time.Duration `env:"IDLE_TIMEOUT"`
}

type Token struct {
	Token string `json:"token"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	UserRole  `json:"role"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
}

var DummyUserUUID uuid.UUID = uuid.MustParse("750cd68b-5736-4abf-a736-578a428d6c01")
var DummyAdminUUID uuid.UUID = uuid.MustParse("05191ec5-abca-41c5-88bb-16641146fe99")
var SecretKey []byte

type UserRole struct {
	Role string `json:"role"`
}

func (ur UserRole) Valid() bool {
	switch ur.Role {
	case "admin":
		return true
	case "user":
		return true
	default:
		return false
	}
}

type Database interface {
	CreateSchedule(schedule *Schedule, roomID uuid.UUID, user User) error
	CreateSlot(slot Slot, startSlot time.Time, endSlot time.Time) error
	GetSlotsList(roomID uuid.UUID, dateRequest time.Time) ([]Slot, error)
	GetSlotBySlotID(slotID uuid.UUID) (Slot, error)
	GetScheduleByRoomID(roomID uuid.UUID) (Schedule, error)

	CreateBooking(booking *Booking) error
	GetAdminBookingList(pagination *Pagination, user User) ([]Booking, error)
	GetUserBookingList(user User) ([]Booking, error)
	GetBookingByBookingID(bookingID uuid.UUID) (Booking, error)
	GetBookingBySlotID(slotID uuid.UUID) (Booking, error)
	CancelBookingByBookingID(bookingID uuid.UUID) (Booking, error)

	CreateRoom(serID uuid.UUID, room *Room) error
	GetRoom(roomID uuid.UUID) (Room, error)
	GetRoomsList() ([]Room, error)
}
type Room struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	Capacity    int        `json:"capacity,omitempty"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
}

type Schedule struct {
	ID         uuid.UUID `json:"id,omitempty"`
	RoomID     uuid.UUID `json:"roomId"`
	DaysOfWeek []int     `json:"daysOfWeek"`
	StartTime  string    `json:"startTime"`
	EndTime    string    `json:"endTime"`
}

type Slot struct {
	ID     uuid.UUID `json:"id"`
	RoomID uuid.UUID `json:"roomId"`
	Start  string    `json:"start"`
	End    string    `json:"end"`
}

type Booking struct {
	ID             uuid.UUID     `json:"id"`
	SlotID         uuid.UUID     `json:"slotId"`
	UserID         uuid.UUID     `json:"userId"`
	Status         BookingStatus `json:"status"`
	ConferenceLink string        `json:"conferenceLink,omitempty"`
	CreatedAt      *time.Time    `json:"createdAt,omitempty"`
}

type BookingStatus string

func (bs BookingStatus) Valid() bool {
	switch bs {
	case "active":
		return true
	case "cancelled":
		return true
	default:
		return false
	}
}

type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
}

type ErrorResponse string

const (
	INVALID_REQUEST     ErrorResponse = "invalid request"
	UNAUTHORIZED        ErrorResponse = "unauthorized"
	NOT_FOUND           ErrorResponse = "not found"
	ROOM_NOT_FOUND      ErrorResponse = "room not found"
	SLOT_NOT_FOUND      ErrorResponse = "slot not found"
	SLOT_ALREADY_BOOKED ErrorResponse = "slot already booked"
	BOOKING_NOT_FOUND   ErrorResponse = "booking not found"
	FORBIDDEN           ErrorResponse = "forbidden"
	SCHEDULE_EXISTS     ErrorResponse = "schedule exists"
)

type InternalErrorResponse string

const INTERNAL_ERROR InternalErrorResponse = "internal_error"
