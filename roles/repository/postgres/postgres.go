package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/rafaceo/go-test-auth/roles/domain"
	"github.com/rafaceo/go-test-auth/roles/repository"
)

type roleRepo struct {
	db *sqlx.DB
}

func NewRoleRepository(db *sqlx.DB) repository.RoleRepository {
	return &roleRepo{db: db}
}

func (r *roleRepo) AddRole(ctx context.Context, roleName, roleNameRu, notes string, rights map[string][]string) error {
	query := `INSERT INTO roles (role_name, role_name_ru, notes, rights) 
	          VALUES ($1, $2, $3, $4)`

	rightsJSON, err := json.Marshal(rights)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, roleName, roleNameRu, notes, rightsJSON)
	return err
}

func (r *roleRepo) EditRole(ctx context.Context, roleID int, roleName, roleNameRu, notes string, rights map[string][]string) error {
	query := `UPDATE roles 
	          SET role_name = $1, role_name_ru = $2, notes = $3, rights = $4
	          WHERE role_id = $5`

	rightsJSON, err := json.Marshal(rights)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, roleName, roleNameRu, notes, rightsJSON, roleID)
	return err
}

func (r *roleRepo) GetRoles(ctx context.Context) ([]domain.Role, error) {
	query := `SELECT role_id, role_name, role_name_ru, notes, rights FROM roles`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []domain.Role
	for rows.Next() {
		var role domain.Role
		var rightsJSON []byte

		if err := rows.Scan(&role.ID, &role.Name, &role.NameRu, &role.Notes, &rightsJSON); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(rightsJSON, &role.Rights); err != nil {
			return nil, err
		}

		roles = append(roles, role)
	}

	return roles, nil
}

func (r *roleRepo) GetRoleRights(ctx context.Context, roleID int) (map[string][]string, error) {
	query := `SELECT rights FROM roles WHERE role_id = $1`

	var rightsJSON []byte
	err := r.db.QueryRowContext(ctx, query, roleID).Scan(&rightsJSON)
	if err != nil {
		return nil, err
	}

	var rights map[string][]string
	if err := json.Unmarshal(rightsJSON, &rights); err != nil {
		return nil, err
	}

	return rights, nil
}

func (r *roleRepo) DeleteRole(ctx context.Context, roleID int) error {
	deleteQuery := `DELETE FROM roles WHERE role_id = $1`
	result, err := r.db.ExecContext(ctx, deleteQuery, roleID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("роль не найдена")
	}

	return nil
}

func (r *roleRepo) AssignRoleToUser(ctx context.Context, userID int, roleID int, merge bool) error {
	userRights, err := r.GetRoleRights(ctx, userID)
	if err != nil {
		return err
	}

	roleRights, err := r.GetRoleRights(ctx, roleID)
	if err != nil {
		return err
	}

	var roleName, roleNameRu, notes string
	err = r.db.QueryRowContext(ctx, "SELECT role_name, role_name_ru, notes FROM roles WHERE role_id = $1", userID).
		Scan(&roleName, &roleNameRu, &notes)
	if err != nil {
		return err
	}

	if merge {
		for key, roleValues := range roleRights {
			userValues, exists := userRights[key]
			if exists {
				rightsSet := make(map[string]struct{})
				for _, v := range userValues {
					rightsSet[v] = struct{}{}
				}
				for _, v := range roleValues {
					rightsSet[v] = struct{}{}
				}

				mergedValues := make([]string, 0, len(rightsSet))
				for v := range rightsSet {
					mergedValues = append(mergedValues, v)
				}
				userRights[key] = mergedValues
			} else {
				userRights[key] = roleValues
			}
		}
	} else {
		userRights = roleRights
	}

	return r.EditRole(ctx, userID, roleName, roleNameRu, notes, userRights)
}
