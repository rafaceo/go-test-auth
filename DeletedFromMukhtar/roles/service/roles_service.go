package service

import (
	"context"

	"github.com/google/uuid"
	"gitlab.fortebank.com/forte-market/apps/user-profile-api/src/roles/domain"
	"gitlab.fortebank.com/forte-market/apps/user-profile-api/src/roles/repository"
)

type RoleService interface {
	AddRole(ctx context.Context, role domain.Role) (uuid.UUID, error)
	GetRoleByID(ctx context.Context, roleID uuid.UUID) (*domain.Role, error)
	GetAllRoles(ctx context.Context) ([]domain.Role, error)
	UpdateRole(ctx context.Context, role domain.Role) error
	DeleteRole(ctx context.Context, roleID uuid.UUID) error
}

type roleService struct {
	repo repository.RoleRepository
}

// NewRoleService создаёт сервис управления ролями
func NewRoleService(repo repository.RoleRepository) RoleService {
	return &roleService{repo: repo}
}

func (s *roleService) AddRole(ctx context.Context, role domain.Role) (uuid.UUID, error) {
	return s.repo.CreateRole(ctx, role)
}

func (s *roleService) GetRoleByID(ctx context.Context, roleID uuid.UUID) (*domain.Role, error) {
	return s.repo.GetRoleByID(ctx, roleID)
}

func (s *roleService) GetAllRoles(ctx context.Context) ([]domain.Role, error) {
	return s.repo.GetAllRoles(ctx)
}

func (s *roleService) UpdateRole(ctx context.Context, role domain.Role) error {
	return s.repo.UpdateRole(ctx, role)
}

func (s *roleService) DeleteRole(ctx context.Context, roleID uuid.UUID) error {
	return s.repo.DeleteRole(ctx, roleID)
}
