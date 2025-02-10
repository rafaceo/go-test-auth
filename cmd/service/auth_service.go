package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"

	//"gitlab.fortebank.com/forte-market/apps/user-profile-api/src/authenticate/domain"
	//"gitlab.fortebank.com/forte-market/apps/user-profile-api/src/authenticate/repository"
	"github.com/rafaceo/go-test-auth/cmd/domain"
	"github.com/rafaceo/go-test-auth/cmd/repository"
)

type AuthService interface {
	Login(ctx context.Context, phone, password string) (string, string, error)
	RefreshToken(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}

type authService struct {
	repo repository.AuthRepository
}

func NewAuthService(repo repository.AuthRepository) AuthService {
	return &authService{repo: repo}
}

func (s *authService) Login(ctx context.Context, phone, password string) (string, string, error) {
	if phone == "" || password == "" {
		return "", "", errors.New("missing required fields") // Будем возвращать 400
	}

	user, err := s.repo.GetUserByPhone(ctx, phone)
	if err != nil {
		return "", "", errors.New("invalid credentials") // Будем возвращать 401
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", "", errors.New("invalid credentials") // Будем возвращать 401
	}

	// Генерация access_token
	accessToken, err := domain.GenerateAccessToken(user.ID, user.Phone, user.Roles, user.Entitlements)
	if err != nil {
		return "", "", errors.New("failed to generate access token") // Будем возвращать 500
	}

	// Генерация refresh_token
	refreshToken := domain.GenerateRefreshToken()

	// Сохранение refresh_token в базе
	err = s.repo.SaveRefreshToken(ctx, user.ID.String(), refreshToken)
	if err != nil {
		fmt.Println("Ошибка сохранения refresh-токена:", err)
		return "", "", errors.New("failed to save refresh token") // Будем возвращать 500
	}

	return accessToken, refreshToken, nil
}

func (s *authService) RefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var req domain.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.RefreshToken == "" {
		http.Error(w, "Некорректный запрос", http.StatusBadRequest)
		return
	}

	// Проверяем, есть ли этот refresh_token в БД
	userID, err := s.repo.GetUserIDByRefreshToken(context.Background(), req.RefreshToken)
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
	accessToken, err := newAccessToken.SignedString(domain.SecretKey)
	if err != nil {
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	// Генерируем новый refresh_token (UUID) и обновляем в БД
	newRefreshToken := uuid.NewString()
	if err := s.repo.UpdateRefreshToken(context.Background(), userID, newRefreshToken); err != nil {
		http.Error(w, "Ошибка обновления токена", http.StatusInternalServerError)
		return
	}

	// Отправляем новые токены клиенту
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(domain.RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	})
}

func (s *authService) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Убедимся, что заголовок Content-Type = application/json
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Неверный Content-Type", http.StatusBadRequest)
		return
	}

	var req domain.LogoutRequest

	// Декодируем JSON-запрос
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.RefreshToken == "" {
		http.Error(w, "Некорректный запрос", http.StatusBadRequest)
		return
	}

	// Вызываем метод репозитория для удаления токена
	err = s.repo.RevokeToken(context.Background(), req.RefreshToken)
	if err != nil {
		log.Println("Ошибка при удалении refreshToken:", err)
		http.Error(w, "Токен не найден", http.StatusForbidden)
		return
	}

	// Отправляем успешный ответ
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Выход выполнен успешно"}`))
}
