package e2e_tests

import (
	"avito_tech_backend/internal/config"
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCancelEndToEnd(t *testing.T) {
	baseURL := "http://localhost:8080"

	var user config.User
	user.Role = "user"

	userToken, err := getAuthToken(baseURL, user)

	if err != nil {
		slog.Error(err.Error())
	}
	assert.NoError(t, err)
	assert.NotZero(t, userToken)

	bookings, err := getUserBookings(baseURL, userToken)

	if len(bookings) == 0 {
		slog.Error("Bookings is empty")
		return
	}

	booking := bookings[0]

	booking, err = cancelBooking(baseURL, userToken, booking)
	assert.NoError(t, err)
	assert.True(t, booking.Status == "cancelled")

}

func cancelBooking(baseURL string, token config.Token, booking config.Booking) (config.Booking, error) {
	client := &http.Client{}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/bookings/%s/cancel", baseURL, booking.ID),
		nil,
	)
	if err != nil {
		slog.Info(err.Error())
		return config.Booking{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token.Token)

	resp, err := client.Do(req)
	if err != nil {
		slog.Info(err.Error())
		return config.Booking{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return config.Booking{}, fmt.Errorf("Cancel booking failed with status: %d", resp.StatusCode)
	}

	var result config.Booking
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		slog.Info(err.Error())
		return config.Booking{}, err
	}

	return result, nil
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

func getUserBookings(baseURL string, token config.Token) ([]config.Booking, error) {
	client := &http.Client{}

	req, err := http.NewRequest(
		"GET",
		baseURL+"/bookings/my",
		nil,
	)

	if err != nil {
		slog.Info(err.Error())
		return []config.Booking{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token.Token)

	resp, err := client.Do(req)
	if err != nil {
		return []config.Booking{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []config.Booking{}, fmt.Errorf("booking failed with status: %d", resp.StatusCode)
	}

	var result []config.Booking
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return []config.Booking{}, err
	}

	return result, nil
}
