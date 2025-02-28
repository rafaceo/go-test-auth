package postgres

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rafaceo/go-test-auth/user_contexts/domain"
	"github.com/rafaceo/go-test-auth/user_contexts/repository"
)

type userContextRepo struct {
	db *sqlx.DB
}

func NewUserContextRepository(db *sqlx.DB) repository.UserContextRepository {
	return &userContextRepo{db: db}
}

func (r *userContextRepo) AddUserContext(ctx context.Context, userCtx domain.UserContext) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO users_contexts (user_id, merchant_id, global) 
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, merchant_id) 
		DO UPDATE SET global = EXCLUDED.global
	`, userCtx.UserID, userCtx.MerchantID, userCtx.Global)
	return err
}

func (r *userContextRepo) EditUserContext(ctx context.Context, userID uuid.UUID, global bool) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE users_contexts 
		SET global = $2
		WHERE user_id = $1
	`, userID, global)
	return err
}

func (r *userContextRepo) GetUserContexts(ctx context.Context, userID uuid.UUID) (domain.UserContext, error) {
	var userCtx domain.UserContext

	err := r.db.QueryRowContext(ctx, `
		SELECT user_id, COALESCE(string_agg(merchant_id, ','), ''), global
		FROM users_contexts
		WHERE user_id = $1
		GROUP BY user_id, global
	`, userID).Scan(&userCtx.UserID, &userCtx.MerchantID, &userCtx.Global)

	if err != nil {
		return domain.UserContext{}, err
	}

	return userCtx, nil
}

func (r *userContextRepo) DeleteUserContext(ctx context.Context, userID uuid.UUID, merchantID string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE users_contexts 
		SET merchant_id = CASE 
			WHEN merchant_id = $2 THEN '' 
			ELSE merchant_id 
		END
		WHERE user_id = $1 AND merchant_id = $2
	`, userID, merchantID)

	return err
}

func (r *userContextRepo) DeleteAllUserContexts(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE users_contexts 
		SET merchant_id = '' 
		WHERE user_id = $1
	`, userID)

	return err
}
