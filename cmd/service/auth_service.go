package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/rafaceo/go-test-auth/cmd/domain"
	"github.com/rafaceo/go-test-auth/cmd/repository"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type AuthService interface {
	Login(ctx context.Context, phone, password string) (string, string, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
	Logout(ctx context.Context, refreshToken string) error
	Register(ctx context.Context, phone string, email string, password string, firstName string, lastName string) (string, error)
}

type authService struct {
	repo      repository.AuthRepository
	jwtSecret string
}

func NewAuthService(repo repository.AuthRepository, jwtSecret string) AuthService {
	return &authService{repo: repo, jwtSecret: jwtSecret}
}

func (s *authService) Register(ctx context.Context, phone string, email, password, firstName, lastName string) (string, error) {
	exists, err := s.repo.UserExists(ctx, phone, email)
	if err != nil {
		return "", err
	}
	if exists {
		return "", errors.New("пользователь уже существует")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	err = s.repo.CreateUser(ctx, phone, email, string(hashedPassword), firstName, lastName)
	if err != nil {
		return "", err
	}

	return "Пользователь успешно зарегистрирован", nil
}

func (s *authService) Login(ctx context.Context, phone, password string) (string, string, error) {
	user, err := s.repo.GetUserByPhone(ctx, phone)
	if err != nil {
		return "", "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	// Генерация access_token
	accessToken, err := domain.GenerateAccessToken(user.ID, user.Phone, user.Roles, user.Entitlements, s.jwtSecret)
	if err != nil {
		return "", "", errors.New("failed to generate access token")
	}

	// Генерация refresh_token
	refreshToken := domain.GenerateRefreshToken()

	// Сохранение refresh_token в базе
	err = s.repo.SaveRefreshToken(ctx, user.ID.String(), refreshToken)
	if err != nil {
		fmt.Println("Ошибка сохранения refresh-токена:", err)
		return "", "", errors.New("failed to save refresh token")
	}

	return accessToken, refreshToken, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	userID, err := s.repo.GetUserIDByRefreshToken(ctx, refreshToken)
	if err != nil {
		log.Printf("Error getting user ID by refresh token: %v", err)
		return "", "", errors.New("invalid refresh token")
	}

	expirationTime := time.Now().Add(15 * time.Minute).Unix()
	newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userID,
		"exp": expirationTime,
	})
	accessToken, err := newAccessToken.SignedString(domain.SecretKey)
	if err != nil {
		return "", "", errors.New("failed to generate access token")
	}

	newRefreshToken := uuid.NewString()
	if err := s.repo.UpdateRefreshToken(ctx, userID, newRefreshToken); err != nil {
		return "", "", errors.New("failed to update refresh token")
	}

	return accessToken, newRefreshToken, nil
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	_, err := s.repo.GetUserIDByRefreshToken(ctx, refreshToken)
	if err != nil {
		log.Printf("Error getting user ID by refresh token: %v", err)
		return errors.New("invalid refresh token")
	}
	return s.repo.DeleteRefreshToken(ctx, refreshToken) // Передаём токен, а не userID

}
