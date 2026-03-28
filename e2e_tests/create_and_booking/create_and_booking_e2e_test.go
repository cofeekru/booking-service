package e2e_tests

import (
	"avito_tech_backend/internal/config"
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEndToEndFlow(t *testing.T) {
	baseURL := "http://localhost:8080"

	var admin config.User
	admin.Role = "admin"

	adminToken, err := getAuthToken(baseURL, admin)

	if err != nil {
		slog.Error(err.Error())
	}
	assert.NoError(t, err)
	assert.NotZero(t, adminToken)

	var room config.Room
	room.Name = "MyBestRoom"
	room.Description = "This is my test and best room"
	room.Capacity = 100

	room, err = createRoom(baseURL, adminToken, room)

	if err != nil {
		slog.Error(err.Error())
	}
	assert.NoError(t, err)
	assert.NotZero(t, room)

	schedule, err := createSchedule(baseURL, adminToken, room)

	assert.NoError(t, err)
	assert.NotZero(t, schedule)

	var user config.User
	user.Role = "user"
	userToken, err := getAuthToken(baseURL, user)
	if err != nil {
		slog.Error(err.Error())
	}
	assert.NoError(t, err)
	assert.NotZero(t, userToken)

	slots, err := getSlotsList(baseURL, userToken, room, time.Now().Add(time.Hour*24).Format(time.DateOnly))

	if len(slots) == 0 {
		slog.Error("слотов нет")
		return
	}
	assert.NoError(t, err)
	assert.NotZero(t, slots)

	slot := slots[0]

	bookingID, err := bookSlot(baseURL, userToken, slot)
	if err != nil {
		slog.Error("no free slots")
	}
	assert.NoError(t, err)
	assert.NotZero(t, bookingID)

}

func getAuthToken(baseURL string, user config.User) (config.Token, error) {
	client := &http.Client{}

	reqBody, _ := json.Marshal(config.UserRole{Role: user.Role})
	req, err := http.NewRequest(
		"POST",
		baseURL+"/dummyLogin",
		bytes.NewBuffer(reqBody))
	if err != nil {
		slog.Info(err.Error())
		return config.Token{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		slog.Info(err.Error())
		return config.Token{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return config.Token{}, fmt.Errorf("login failed with status: %d", resp.StatusCode)
	}

	var token config.Token
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		slog.Info(err.Error())
		return config.Token{}, err
	}

	return token, nil
}

func createRoom(baseURL string, token config.Token, room config.Room) (config.Room, error) {
	client := &http.Client{}

	reqBody, _ := json.Marshal(room)
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/rooms/create", baseURL),
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		slog.Info(err.Error())
		return config.Room{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token.Token)

	resp, err := client.Do(req)
	if err != nil {
		slog.Info(err.Error())
		return config.Room{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return config.Room{}, fmt.Errorf("Room creation failed with status: %d", resp.StatusCode)
	}

	var result config.Room

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		slog.Info(err.Error())
		return config.Room{}, err
	}

	return result, nil
}

func createSchedule(baseURL string, token config.Token, room config.Room) (config.Schedule, error) {
	client := &http.Client{}

	scheduleReq := config.Schedule{
		RoomID:     room.ID,
		DaysOfWeek: []int{1, 2, 3, 4, 5, 6, 7},
		StartTime:  "09:00",
		EndTime:    "21:00",
	}
	reqBody, _ := json.Marshal(scheduleReq)
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/rooms/%s/schedule/create", baseURL, room.ID),
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		slog.Info(err.Error())
		return config.Schedule{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token.Token)

	resp, err := client.Do(req)
	if err != nil {
		slog.Info(err.Error())
		return config.Schedule{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return config.Schedule{}, fmt.Errorf("Schedule creation failed with status: %d", resp.StatusCode)
	}

	var result config.Schedule
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		slog.Info(err.Error())
		return config.Schedule{}, err
	}

	return result, nil
}

func getSlotsList(baseURL string, token config.Token, room config.Room, date string) ([]config.Slot, error) {
	client := &http.Client{}

	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/rooms/%s/slots/list?date=%s", baseURL, room.ID, date),
		nil,
	)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token.Token)
	resp, err := client.Do(req)
	if err != nil {
		slog.Info(err.Error())
		return []config.Slot{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []config.Slot{}, fmt.Errorf("Slots getting failed with status: %d", resp.StatusCode)
	}

	var slots []config.Slot
	if err := json.NewDecoder(resp.Body).Decode(&slots); err != nil {
		slog.Info(err.Error())
		return []config.Slot{}, err
	}

	return slots, nil

}

func bookSlot(baseURL string, token config.Token, slot config.Slot) (config.Booking, error) {
	client := &http.Client{}

	bookingReq := config.Booking{
		SlotID: slot.ID,
	}

	reqBody, _ := json.Marshal(bookingReq)

	req, err := http.NewRequest(
		"POST",
		baseURL+"/bookings/create",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		slog.Info(err.Error())
		return config.Booking{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token.Token)

	resp, err := client.Do(req)
	if err != nil {
		return config.Booking{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return config.Booking{}, fmt.Errorf("Booking failed with status: %d", resp.StatusCode)
	}

	var result config.Booking
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return config.Booking{}, err
	}

	return result, nil
}
