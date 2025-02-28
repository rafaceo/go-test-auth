package repository

import (
	"context"
	"github.com/rafaceo/go-test-auth/rights/domain"
)

type RightsRepository interface {
	AddRights(ctx context.Context, module string, action []string) error
	EditRight(ctx context.Context, id string, module string, action []string) error
	GetAllRights(ctx context.Context) ([]domain.Right, error)
	GetRightByName(ctx context.Context, module string) (*domain.Right, error)
	GetRightById(ctx context.Context, id string) (*domain.Right, error)
	DeleteRight(ctx context.Context, id string) error
}
