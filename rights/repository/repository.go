package repository

import (
	"context"
	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, phone string, passwordHash string) error
	EditUser(ctx context.Context, id uuid.UUID, phone, password string) error
	GrantRightsToUser(ctx context.Context, id uuid.UUID, rights map[string][]string) error
	EditRightsToUser(ctx context.Context, id uuid.UUID, rights map[string][]string) error
	RevokeRightsFromUser(ctx context.Context, id uuid.UUID, rights map[string][]string) error
	GetUser(ctx context.Context, id uuid.UUID) (string, string, string, string, error)
	GetUserRights(ctx context.Context, id uuid.UUID) (map[string][]string, error)
}
