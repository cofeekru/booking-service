package services_test

import (
	"testing"
	"time"

	"avito_tech_backend/internal/config"
	"avito_tech_backend/internal/services"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestDummyLogin_AdminRole(t *testing.T) {
	userRole := config.UserRole{Role: "admin"}

	token, err := services.DummyLogin(userRole)

	if err != nil {
		t.Fatalf("DummyLogin() for admin role returned error: %v", err)
	}

	if token == "" {
		t.Error("DummyLogin() returned empty token for admin role")
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.SecretKey), nil
	})
	if err != nil {
		t.Fatalf("Failed to parse generated JWT: %v", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("Failed to get claims from JWT")
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		t.Error("user_id claim not found or not a string")
	} else if userIDStr != config.DummyAdminUUID.String() {
		t.Errorf("user_id claim = %s, want %s", userIDStr, config.DummyAdminUUID.String())
	}

	role, ok := claims["role"].(string)
	if !ok {
		t.Error("role claim not found or not a string")
	} else if role != "admin" {
		t.Errorf("role claim = %s, want 'admin'", role)
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		t.Error("exp claim not found or not a number")
	} else {
		expTime := time.Unix(int64(exp), 0)
		expectedExp := time.Now().Add(time.Hour * 1)
		if expTime.Before(expectedExp.Add(-time.Minute)) || expTime.After(expectedExp.Add(time.Minute)) {
			t.Errorf("exp claim = %v, want approximately %v", expTime, expectedExp)
		}
	}
}

func TestDummyLogin_UserRole(t *testing.T) {
	userRole := config.UserRole{Role: "user"}

	token, err := services.DummyLogin(userRole)

	if err != nil {
		t.Fatalf("DummyLogin() for user role returned error: %v", err)
	}

	if token == "" {
		t.Error("DummyLogin() returned empty token for user role")
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.SecretKey), nil
	})
	if err != nil {
		t.Fatalf("Failed to parse generated JWT: %v", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("Failed to get claims from JWT")
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		t.Error("user_id claim not found or not a string")
	} else if userIDStr != config.DummyUserUUID.String() {
		t.Errorf("user_id claim = %s, want %s", userIDStr, config.DummyUserUUID.String())
	}

	role, ok := claims["role"].(string)
	if !ok {
		t.Error("role claim not found or not a string")
	} else if role != "user" {
		t.Errorf("role claim = %s, want 'user'", role)
	}
}

func TestDummyLogin_UnknownRole(t *testing.T) {
	unknownRole := config.UserRole{Role: "unknown"}

	token, err := services.DummyLogin(unknownRole)

	if err != nil {
		t.Fatalf("DummyLogin() for unknown role returned error: %v", err)
	}

	if token == "" {
		t.Error("DummyLogin() returned empty token for unknown role")
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.SecretKey), nil
	})
	if err != nil {
		t.Fatalf("Failed to parse generated JWT: %v", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("Failed to get claims from JWT")
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		t.Error("user_id claim not found or not a string")
	} else if userUUID, _ := uuid.Parse(userIDStr); userUUID != uuid.Nil {
		t.Errorf("user_id claim = %s, want empty UUID", userIDStr)
	}

	role, ok := claims["role"].(string)
	if !ok {
		t.Error("role claim not found or not a string")
	} else if role != "unknown" {
		t.Errorf("role claim = %s, want 'unknown'", role)
	}
}
