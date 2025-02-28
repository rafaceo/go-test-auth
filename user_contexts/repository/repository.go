package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/rafaceo/go-test-auth/user_contexts/domain"
)

type UserContextRepository interface {
	AddUserContext(ctx context.Context, userCtx domain.UserContext) error
	EditUserContext(ctx context.Context, userID uuid.UUID, global bool) error
	GetUserContexts(ctx context.Context, userID uuid.UUID) (domain.UserContext, error)
	DeleteUserContext(ctx context.Context, userID uuid.UUID, merchantID string) error
	DeleteAllUserContexts(ctx context.Context, userID uuid.UUID) error
}
