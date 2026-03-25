package handlers

import (
	"avito_tech_backend/internal/config"
	"avito_tech_backend/internal/services"
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func DummyLoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var userRole config.UserRole
		err := json.NewDecoder(r.Body).Decode(&userRole)

		if err != nil || !userRole.Valid() {
			log.Println("Неверный запрос: недопустимое значение роли или тела запроса ", err, userRole.Role)
			http.Error(w, string(config.INVALID_REQUEST), http.StatusBadRequest)
			return
		}

		token, err := services.DummyLogin(userRole)

		if err != nil {
			log.Println("Ошибка создания JWT токена ", err)
			http.Error(w, string(config.INTERNAL_ERROR), http.StatusInternalServerError)
			return
		}

		var jwtToken config.Token = config.Token{
			Token: token,
		}

		w.Header().Set("Authorization", "Bearer "+jwtToken.Token)
		err = json.NewEncoder(w).Encode(jwtToken)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, string(config.INTERNAL_ERROR), http.StatusInternalServerError)
			return
		}
	}
}

func AuthMiddlewareHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, string(config.UNAUTHORIZED), http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			return config.SecretKey, nil
		})

		if err != nil || !token.Valid {
			slog.Error("Не авторизован")
			http.Error(w, string(config.UNAUTHORIZED), http.StatusUnauthorized)
			return
		}

		userInfo := token.Claims.(jwt.MapClaims)

		userID, _ := userInfo["user_id"].(string)
		userRole, _ := userInfo["role"].(string)

		ctx := context.WithValue(r.Context(), "userID", userID)
		ctx = context.WithValue(ctx, "userRole", userRole)

		next.ServeHTTP(w, r.WithContext(ctx))
	})

}
