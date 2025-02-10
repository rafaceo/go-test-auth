package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/rafaceo/go-test-auth/cmd/domain"
	repo "github.com/rafaceo/go-test-auth/cmd/repository"
)

type authRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) repo.AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) CreateUser(ctx context.Context, phone string, email, password_hash, firstName, lastName string) error {
	query := `INSERT INTO users_profiles (phone, email, password_hash, first_name, last_name) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, phone, email, password_hash, firstName, lastName)
	return err
}

func (r *authRepository) UserExists(ctx context.Context, phone string, email string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM users_profiles WHERE phone = $1 OR email = $2`
	err := r.db.GetContext(ctx, &count, query, phone, email)
	if err != nil {
		return false, err
	}
	return count > 0, nil
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

func (r *authRepository) DeleteRefreshToken(ctx context.Context, refreshToken string) error {
	query := `UPDATE users_profiles SET refresh_token = NULL WHERE refresh_token = $1`
	_, err := r.db.ExecContext(ctx, query, refreshToken)
	return err
}
