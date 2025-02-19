package repository

import (
	"context"
)

type UserRepository interface {
	CreateUser(ctx context.Context, name string) error
	EditUser(ctx context.Context, id uint, name, context string) error
	GrantRightsToUser(ctx context.Context, id uint, rights map[string][]string) error
	EditRightsToUser(ctx context.Context, id uint, rights map[string][]string) error
	RevokeRightsFromUser(ctx context.Context, id uint, rights map[string][]string) error
	GetUser(ctx context.Context, id uint) (string, string, map[string][]string, error)
}
