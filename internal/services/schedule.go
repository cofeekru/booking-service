package services

import (
	"avito_tech_backend/internal/config"

	"time"

	"github.com/google/uuid"
)

func ScheduleExist(storage config.Database, schedule *config.Schedule) bool {
	_, err := storage.GetScheduleByRoomID(schedule.RoomID)
	return err == nil
}

func ValidSchedule(schedule config.Schedule) bool {
	daysOfWeek := schedule.DaysOfWeek

	if len(daysOfWeek) > 7 || len(daysOfWeek) < 1 {
		return false
	}

	var uniqueValue map[int]bool = make(map[int]bool)

	for _, value := range daysOfWeek {
		if value < 1 || value > 7 {
			return false
		}
		uniqueValue[value] = true
	}

	return len(uniqueValue) == len(daysOfWeek)
}

func ScheduleCreate(storage config.Database, schedule *config.Schedule, roomID uuid.UUID, user config.User) error {
	newUUID, err := uuid.NewV6()

	if err != nil {
		return err
	}
	schedule.ID = newUUID
	err = storage.CreateSchedule(schedule, roomID, user)

	if err != nil {
		return err
	}

	startTime, _ := time.Parse("15:04", schedule.StartTime)
	endTime, _ := time.Parse("15:04", schedule.EndTime)

	if endTime.Sub(startTime).Minutes() < 30 {
		return nil
	}

	for _, weekDay := range schedule.DaysOfWeek {
		today := time.Now().Weekday()

		var startSlot, endSlot time.Time
		var diffDay int
		if int(today) < weekDay && today != time.Sunday {
			diffDay = weekDay - int(today)
		} else if today != time.Sunday {
			diffDay = -(int(today) - weekDay)
		} else {
			diffDay = -(7 - weekDay)
		}
		startSlot = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()+diffDay, startTime.Hour(), startTime.Minute(), 0, 0, time.Local)
		endSlot = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()+diffDay, endTime.Hour(), endTime.Minute(), 0, 0, time.Local)

		if diffDay < 0 {
			startSlot = startSlot.Add(time.Hour * 168)
			endSlot = endSlot.Add(time.Hour * 168)
		}

		countWeek := 5
		for range countWeek {
			intermediateStartSlot := startSlot

			for endSlot.Sub(intermediateStartSlot).Minutes() > 30 {
				var slot config.Slot
				slot.ID = uuid.New()
				slot.RoomID = schedule.RoomID

				err := storage.CreateSlot(slot, intermediateStartSlot, intermediateStartSlot.Add(time.Minute*30))
				if err != nil {
					return err
				}
				intermediateStartSlot = intermediateStartSlot.Add(time.Minute * 40)
			}
			startSlot = startSlot.Add(time.Hour * 168)
			endSlot = endSlot.Add(time.Hour * 168)
		}
	}

	return nil
}
