package repository

import (
	"context"
	"github.com/rafaceo/go-test-auth/roles/domain"
)

type RoleRepository interface {
	AddRole(ctx context.Context, roleName, roleNameRu, notes string, rights map[string][]string) error
	EditRole(ctx context.Context, roleID int, roleName, roleNameRu, notes string, rights map[string][]string) error
	GetRoles(ctx context.Context) ([]domain.Role, error)
	GetRoleRights(ctx context.Context, roleID int) (map[string][]string, error)
	DeleteRole(ctx context.Context, roleID int) error
	AssignRoleToUser(ctx context.Context, userID int, roleID int, merge bool) error
}
