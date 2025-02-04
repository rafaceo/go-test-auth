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
