package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/jmoiron/sqlx"
	"gitlab.fortebank.com/forte-market/apps/user-profile-api/src/roles/repository"
	"time"

	"github.com/google/uuid"
	"gitlab.fortebank.com/forte-market/apps/user-profile-api/src/roles/domain"
)

type roleRepo struct {
	db *sqlx.DB
}

// NewRoleRepository создаёт новый репозиторий ролей
func NewRoleRepository(db *sqlx.DB) repository.RoleRepository {
	return &roleRepo{db: db}
}

// CreateRole добавляет новую роль
func (r *roleRepo) CreateRole(ctx context.Context, role domain.Role) (uuid.UUID, error) {
	query := `INSERT INTO roles (role_id, role_name, rights, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING role_id`
	rightsJSON, _ := json.Marshal(role.Rights)
	roleID := uuid.New()

	err := r.db.QueryRowContext(ctx, query, roleID, role.RoleName, rightsJSON, time.Now(), time.Now()).Scan(&roleID)
	if err != nil {
		return uuid.Nil, err
	}
	return roleID, nil
}

// GetRoleByID возвращает роль по ID
func (r *roleRepo) GetRoleByID(ctx context.Context, roleID uuid.UUID) (*domain.Role, error) {
	query := `SELECT role_id, role_name, rights, created_at, updated_at FROM roles WHERE role_id = $1`
	row := r.db.QueryRowContext(ctx, query, roleID)

	var role domain.Role
	var rightsJSON []byte

	if err := row.Scan(&role.RoleID, &role.RoleName, &rightsJSON, &role.CreatedAt, &role.UpdatedAt); err != nil {
		if err != nil {
			return nil, nil
		}
		return nil, err
	}

	json.Unmarshal(rightsJSON, &role.Rights)
	return &role, nil
}

// GetRoleByName возвращает роль по имени
func (r *roleRepo) GetRoleByName(ctx context.Context, roleName string) (*domain.Role, error) {
	query := `SELECT role_id, role_name, rights, created_at, updated_at FROM roles WHERE role_name = $1`
	row := r.db.QueryRowContext(ctx, query, roleName)

	var role domain.Role
	var rightsJSON []byte

	if err := row.Scan(&role.RoleID, &role.RoleName, &rightsJSON, &role.CreatedAt, &role.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	json.Unmarshal(rightsJSON, &role.Rights)
	return &role, nil
}

// GetAllRoles возвращает все роли
func (r *roleRepo) GetAllRoles(ctx context.Context) ([]domain.Role, error) {
	query := `SELECT role_id, role_name, rights, created_at, updated_at FROM roles`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []domain.Role
	for rows.Next() {
		var role domain.Role
		var rightsJSON []byte

		if err := rows.Scan(&role.RoleID, &role.RoleName, &rightsJSON, &role.CreatedAt, &role.UpdatedAt); err != nil {
			return nil, err
		}

		json.Unmarshal(rightsJSON, &role.Rights)
		roles = append(roles, role)
	}
	return roles, nil
}

// UpdateRole обновляет данные роли
func (r *roleRepo) UpdateRole(ctx context.Context, role domain.Role) error {
	query := `UPDATE roles SET role_name = $1, rights = $2, updated_at = $3 WHERE role_id = $4`
	rightsJSON, _ := json.Marshal(role.Rights)

	_, err := r.db.ExecContext(ctx, query, role.RoleName, rightsJSON, time.Now(), role.RoleID)
	return err
}

// DeleteRole удаляет роль
func (r *roleRepo) DeleteRole(ctx context.Context, roleID uuid.UUID) error {
	query := `DELETE FROM roles WHERE role_id = $1`
	_, err := r.db.ExecContext(ctx, query, roleID)
	return err
}
