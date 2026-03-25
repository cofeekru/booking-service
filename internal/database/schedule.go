package database

import (
	"avito_tech_backend/internal/config"
	"encoding/json"

	"github.com/google/uuid"
)

func (storage *Storage) CreateSchedule(schedule *config.Schedule, roomID uuid.UUID, user config.User) error {
	_, err := storage.DB.Exec(`
		INSERT INTO schedule (schedule_id, room_id, days_of_week, start_time, end_time) 
		VALUES ($1, $2, $3, $4, $5)`,
		schedule.ID, roomID, schedule.DaysOfWeek, schedule.StartTime, schedule.EndTime)
	return err
}

func (storage *Storage) GetScheduleByRoomID(roomID uuid.UUID) (config.Schedule, error) {
	var schedule config.Schedule
	var days_of_week []byte

	err := storage.DB.QueryRow(`
		SELECT schedule_id, room_id, days_of_week, start_time, end_time 
		FROM schedule WHERE room_id = $1`, roomID.String()).
		Scan(&schedule.ID, &schedule.RoomID, &days_of_week, &schedule.StartTime, &schedule.EndTime)

	json.Unmarshal(days_of_week, &schedule.DaysOfWeek)

	if err != nil {
		return config.Schedule{}, err
	}

	return schedule, err
}
