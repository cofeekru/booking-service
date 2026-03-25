package e2e_tests

import (
	"avito_tech_backend/internal/config"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type ScheduleRequest struct {
	RoomID    string    `json:"roomId"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

type BookingRequest struct {
	SlotID string `json:"slotId"`
}

func TestEndToEndFlow(t *testing.T) {
	baseURL := "http://localhost:8080"

	adminToken, err := getAuthToken(baseURL, "admin")
	assert.NoError(t, err)
	assert.NotZero(t, adminToken)

	roomID := "test-room-123"
	scheduleID, err := createSchedule(baseURL, adminToken, roomID)
	assert.NoError(t, err)
	assert.NotZero(t, scheduleID)

	userToken, err := getAuthToken(baseURL, "user")
	assert.NoError(t, err)
	assert.NotZero(t, userToken)

	bookingID, err := bookSlot(baseURL, userToken, scheduleID)
	assert.NoError(t, err)
	assert.NotZero(t, bookingID)

	t.Logf("E2E тест успешно пройден! Booking ID: %s", bookingID)
}

func getAuthToken(baseURL, role string) (string, error) {
	client := &http.Client{}

	reqBody, _ := json.Marshal(config.UserRole{Role: role})
	req, err := http.NewRequest("POST", baseURL+"/dummyLogin", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("login failed with status: %d", resp.StatusCode)
	}

	var token config.Token
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return "", err
	}

	return token.Token, nil
}

func createSchedule(baseURL, token, roomID string) (string, error) {
	client := &http.Client{}

	now := time.Now()
	scheduleReq := ScheduleRequest{
		RoomID:    roomID,
		StartTime: now.Add(24 * time.Hour),
		EndTime:   now.Add(25 * time.Hour),
	}

	reqBody, _ := json.Marshal(scheduleReq)
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/rooms/%s/schedule/create", baseURL, roomID),
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("schedule creation failed with status: %d", resp.StatusCode)
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result["id"], nil
}

func bookSlot(baseURL, token, scheduleID string) (string, error) {
	client := &http.Client{}

	bookingReq := BookingRequest{SlotID: scheduleID}
	reqBody, _ := json.Marshal(bookingReq)

	req, err := http.NewRequest(
		"POST",
		baseURL+"/bookings/create",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("booking failed with status: %d", resp.StatusCode)
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result["id"], nil
}
