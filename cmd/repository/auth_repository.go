package repository

import (
	"context"
	"github.com/jmoiron/sqlx"

	//"gitlab.fortebank.com/forte-market/apps/user-profile-api/src/authenticate/domain"
	"github.com/rafaceo/go-test-auth/cmd/domain"
)

type AuthRepository interface {
	GetUserByPhone(ctx context.Context, phone string) (*domain.UserProfile, error)
	SaveRefreshToken(ctx context.Context, userID string, refreshToken string) error
	GetUserIDByRefreshToken(ctx context.Context, refreshToken string) (string, error)    // Новый метод
	UpdateRefreshToken(ctx context.Context, userID string, newRefreshToken string) error // Новый метод
}

type authRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) GetUserByPhone(ctx context.Context, phone string) (*domain.UserProfile, error) {
	var user domain.UserProfile
	query := `SELECT id, phone, password_hash FROM users_profiles WHERE phone = $1`
	err := r.db.QueryRowContext(ctx, query, phone).Scan(&user.ID, &user.Phone, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) SaveRefreshToken(ctx context.Context, userID string, refreshToken string) error {
	query := `UPDATE users_profiles SET refresh_token = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, refreshToken, userID)
	return err
}

// Новый метод для получения userID по refresh_token.
func (r *authRepository) GetUserIDByRefreshToken(ctx context.Context, refreshToken string) (string, error) {
	var userID string
	query := `SELECT id FROM users_profiles WHERE refresh_token = $1`
	err := r.db.QueryRowContext(ctx, query, refreshToken).Scan(&userID)
	if err != nil {
		return "", err
	}
	return userID, nil
}

// Новый метод для обновления refresh_token.
func (r *authRepository) UpdateRefreshToken(ctx context.Context, userID string, newRefreshToken string) error {
	query := `UPDATE users_profiles SET refresh_token = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, newRefreshToken, userID)
	return err
}
