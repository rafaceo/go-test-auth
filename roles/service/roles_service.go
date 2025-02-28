package service

import (
	"context"
	"errors"

	"github.com/rafaceo/go-test-auth/roles/domain"
	"github.com/rafaceo/go-test-auth/roles/repository"
)

type RoleService interface {
	AddRole(ctx context.Context, roleName, roleNameRu, notes string, rights map[string][]string) error
	EditRole(ctx context.Context, roleID int, roleName, roleNameRu, notes string, rights map[string][]string) error
	GetRoles(ctx context.Context) ([]domain.Role, error)
	GetRoleRights(ctx context.Context, roleID int) (map[string][]string, error)
	DeleteRole(ctx context.Context, roleID int) error
	AssignRoleToUser(ctx context.Context, userID int, roleID int, merge bool) error
}

type roleService struct {
	repo repository.RoleRepository
}

// NewRoleService создаёт сервис управления ролями
func NewRoleService(repo repository.RoleRepository) RoleService {
	return &roleService{repo: repo}
}

func (s *roleService) AddRole(ctx context.Context, roleName, roleNameRu, notes string, rights map[string][]string) error {

	s.repo.AddRole(ctx, roleName, roleNameRu, notes, rights)

	return nil
}

func (s *roleService) EditRole(ctx context.Context, roleID int, roleName, roleNameRu, notes string, rights map[string][]string) error {
	if roleID <= 0 {
		return errors.New("invalid role ID")
	}
	if roleName == "" || roleNameRu == "" {
		return errors.New("role name cannot be empty")
	}

	return s.repo.EditRole(ctx, roleID, roleName, roleNameRu, notes, rights)
}

func (s *roleService) GetRoles(ctx context.Context) ([]domain.Role, error) {
	return s.repo.GetRoles(ctx)
}

func (s *roleService) GetRoleRights(ctx context.Context, roleID int) (map[string][]string, error) {
	return s.repo.GetRoleRights(ctx, roleID)
}

func (s *roleService) DeleteRole(ctx context.Context, roleID int) error {
	return s.repo.DeleteRole(ctx, roleID)
}
func (s *roleService) AssignRoleToUser(ctx context.Context, userID int, roleID int, merge bool) error {
	return s.repo.AssignRoleToUser(ctx, userID, roleID, merge)
}
