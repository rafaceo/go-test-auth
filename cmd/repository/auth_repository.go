package repository

import (
	"context"
	"github.com/rafaceo/go-test-auth/cmd/domain"
)

type AuthRepository interface {
	GetUserByPhone(ctx context.Context, phone string) (*domain.UserProfile, error)
	SaveRefreshToken(ctx context.Context, userID string, refreshToken string) error
	GetUserIDByRefreshToken(ctx context.Context, refreshToken string) (string, error) // Новый метод
	UpdateRefreshToken(ctx context.Context, userID string, newRefreshToken string) error
	DeleteRefreshToken(ctx context.Context, refreshToken string) error
	UserExists(ctx context.Context, phone string, email string) (bool, error)
	CreateUser(ctx context.Context, phone string, email, password, firstName, lastName string) error
}
