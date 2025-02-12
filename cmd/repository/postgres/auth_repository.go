package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rafaceo/go-test-auth/cmd/domain"
	repo "github.com/rafaceo/go-test-auth/cmd/repository"
	"time"
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

func (r *authRepository) CreateFiledAttempt(ctx context.Context, phone string) error {
	query := `INSERT INTO failed_logins (phone, attempts, blocked_until) 
              VALUES ($1, 0, NULL) 
              ON CONFLICT (phone) DO NOTHING`

	fmt.Println("TRYING TO INSERT ATTEMPT FOR", phone)
	_, err := r.db.ExecContext(ctx, query, phone)

	if err != nil {
		fmt.Println("SQL ERROR:", err)
	} else {
		fmt.Println("INSERT SUCCESS")
	}

	query = `UPDATE failed_logins 
             SET attempts = attempts + 1 
             WHERE phone = $1 
             RETURNING attempts`

	var attempts int
	err = r.db.QueryRowContext(ctx, query, phone).Scan(&attempts)
	if err != nil {
		fmt.Println("SQL ERROR (increment attempts):", err)
		return err
	}

	var blockDuration = 0 * time.Second // <- Добавил начальное значение
	if attempts == 3 {
		blockDuration = 30 * time.Second
	} else if attempts == 6 {
		blockDuration = 1 * time.Minute
	}

	if blockDuration > 0 {
		blockUntil := time.Now().UTC().Add(blockDuration)
		query = `UPDATE failed_logins 
                 SET blocked_until = $1 
                 WHERE phone = $2`
		_, err = r.db.ExecContext(ctx, query, blockUntil, phone)
		if err != nil {
			fmt.Println("SQL ERROR (block user):", err)
			return err
		}
		fmt.Println("User", phone, "blocked until", blockUntil)
	}

	fmt.Println("INSERT/UPDATE SUCCESS for", phone, "attempts:", attempts)

	return nil
}

func (r *authRepository) CheckBan(ctx context.Context, phone string) error {
	var blockedUntil sql.NullTime
	var attempts int
	fmt.Println("sdsdsdsd")
	query := `SELECT blocked_until, attempts FROM failed_logins WHERE phone = $1`
	err := r.db.QueryRowContext(ctx, query, phone).Scan(&blockedUntil, &attempts)
	fmt.Println("att", attempts, "blocked until", blockedUntil)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		fmt.Println("SQL ERROR (CheckBan):", err)
		return err
	}

	if blockedUntil.Valid && time.Now().UTC().Before(blockedUntil.Time) {
		fmt.Println("CURRENT TIME:", time.Now().UTC())
		fmt.Println("BLOCK UNTIL :", blockedUntil.Time)
		fmt.Println("IS STILL BANNED?", time.Now().UTC().Before(blockedUntil.Time))

		return errors.New("Banned user")
	}

	fmt.Println("attempts:", attempts)
	// 6 attempts
	if attempts >= 6 {
		query = `DELETE FROM failed_logins WHERE phone = $1`
		_, err = r.db.ExecContext(ctx, query, phone)
		if err != nil {
			fmt.Println("SQL ERROR (Delete after ban):", err)
			return err
		}
		fmt.Println("Ban expired, record deleted for", phone)
		return nil
	}

	// 3 attempts
	query = `UPDATE failed_logins SET blocked_until = NULL WHERE phone = $1`
	_, err = r.db.ExecContext(ctx, query, phone)
	if err != nil {
		fmt.Println("SQL ERROR (Unblock after small ban):", err)
		return err
	}

	fmt.Println("Ban expired, attempts kept for", phone)
	return nil
}
