package services

import (
	"log/slog"
	"net/mail"
	"time"

	"avito_tech_backend/internal/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func DummyLogin(userRole config.UserRole) (string, error) {
	var userUUID uuid.UUID
	switch userRole.Role {
	case "admin":
		userUUID = config.DummyAdminUUID
	case "user":
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

func RegisterUser(storage config.Database, user config.User) (config.User, error) {
	hashPassword, err := hashPassword(user.Password)

	if err != nil {
		return config.User{}, nil
	}

	user.ID = uuid.New()
	user.Password = hashPassword
	user.CreatedAt = time.Now().Format(time.RFC3339)

	err = storage.CreateUser(user)
	if err != nil {
		return config.User{}, err
	}
	return user, nil

}

func EmailValidate(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error(err.Error())
		return "", err
	}
	return string(hash), nil
}

func checkPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		slog.Error(err.Error())
	}
	return err == nil
}

func EmailExist(storage config.Database, email string) bool {
	_, err := storage.GetUserByEmail(email)

	return err == nil
}

func UserExist(storage config.Database, userRequest config.User) bool {
	userStorage, err := storage.GetUserByEmail(userRequest.Email)

	if err != nil {
		slog.Error(err.Error())
		return false
	}

	return checkPassword(userStorage.Password, userRequest.Password)
}

func LoginUser(storage config.Database, userInput config.User) (string, error) {
	userInput, _ = storage.GetUserByEmail(userInput.Email)
	claims := jwt.MapClaims{
		"user_id": userInput.ID,
		"role":    userInput.Role,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	}

	jwtToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(config.SecretKey))

	return jwtToken, err
}
