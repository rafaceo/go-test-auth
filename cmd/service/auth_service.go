package service

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	//"gitlab.fortebank.com/forte-market/apps/user-profile-api/src/authenticate/domain"
	//"gitlab.fortebank.com/forte-market/apps/user-profile-api/src/authenticate/repository"
	"github.com/rafaceo/go-test-auth/cmd/domain"
	"github.com/rafaceo/go-test-auth/cmd/repository"
)

type AuthService interface {
	Login(ctx context.Context, phone, password string) (string, string, error)
}

type authService struct {
	repo repository.AuthRepository
}

func NewAuthService(repo repository.AuthRepository) AuthService {
	return &authService{repo: repo}
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
	accessToken, err := domain.GenerateAccessToken(user.ID, user.Phone, user.Roles, user.Entitlements)
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
