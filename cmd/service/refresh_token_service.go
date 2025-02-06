package service

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	cmd "github.com/rafaceo/go-test-auth/cmd/db"
	dom "github.com/rafaceo/go-test-auth/cmd/domain"
)

// RefreshTokenHandler обрабатывает обновление access_token
func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var req dom.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.RefreshToken == "" {
		http.Error(w, "Некорректный запрос", http.StatusBadRequest)
		return
	}

	// Проверяем, есть ли этот refresh_token в БД
	var userID string
	err := cmd.Db.QueryRow("SELECT id FROM users_profiles WHERE refresh_token = $1", req.RefreshToken).Scan(&userID)
	if err != nil {
		http.Error(w, "Недействительный refresh_token", http.StatusForbidden)
		return
	}

	// Генерируем новый access_token (JWT)
	expirationTime := time.Now().Add(15 * time.Minute).Unix()
	newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userID,
		"exp": expirationTime,
	})
	accessToken, err := newAccessToken.SignedString(dom.SecretKey)
	if err != nil {
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	// Генерируем новый refresh_token (UUID) и обновляем в БД
	newRefreshToken := uuid.NewString()
	_, err = cmd.Db.Exec("UPDATE users_profiles SET refresh_token = $1 WHERE id = $2", newRefreshToken, userID)
	if err != nil {
		http.Error(w, "Ошибка обновления токена", http.StatusInternalServerError)
		return
	}

	// Отправляем новые токены клиенту
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dom.RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	})
}
