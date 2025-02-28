package repository

import (
	"context"
	"github.com/google/uuid"
	"gitlab.fortebank.com/forte-market/apps/user-profile-api/src/roles/domain"
)

// RoleRepository интерфейс для работы с ролями
type RoleRepository interface {
	CreateRole(ctx context.Context, role domain.Role) (uuid.UUID, error)
	GetRoleByID(ctx context.Context, roleID uuid.UUID) (*domain.Role, error)
	GetRoleByName(ctx context.Context, roleName string) (*domain.Role, error)
	GetAllRoles(ctx context.Context) ([]domain.Role, error)
	UpdateRole(ctx context.Context, role domain.Role) error
	DeleteRole(ctx context.Context, roleID uuid.UUID) error
}
