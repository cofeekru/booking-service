package services

import (
	"time"

	"avito_tech_backend/internal/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func DummyLogin(userRole config.UserRole) (string, error) {
	var userUUID uuid.UUID
	if userRole.Role == "admin" {
		userUUID = config.DummyAdminUUID
	} else if userRole.Role == "user" {
		userUUID = config.DummyUserUUID
	}

	claims := jwt.MapClaims{
		"user_id": userUUID.String(),
		"role":    userRole.Role,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	}

	jwtToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(config.SecretKey))

	return jwtToken, err
}
