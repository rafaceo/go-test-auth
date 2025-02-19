package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	repo "github.com/rafaceo/go-test-auth/rights/repository"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) repo.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, name string) error {
	query := `INSERT INTO users_rights (name, context, rights) VALUES ($1, null, null)`
	_, err := r.db.ExecContext(ctx, query, name)
	fmt.Println(name, err)
	return err
}

func (r *userRepository) EditUser(ctx context.Context, id uint, name, context string) error {
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM users_rights WHERE id = $1)`
	err := r.db.QueryRowContext(ctx, checkQuery, id).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("user rights with given id not found")
	}

	query := `UPDATE users_rights SET name = $1, context = $2 WHERE id = $3`
	_, err = r.db.ExecContext(ctx, query, name, context, id)
	return err
}
func (r *userRepository) GrantRightsToUser(ctx context.Context, id uint, rights map[string][]string) error {
	rightsJSON, err := json.Marshal(rights)
	if err != nil {
		return err
	}

	// Запрос обновления прав, если текущие права NULL или пустые
	query := `
		UPDATE users_rights 
		SET rights = $1 
		WHERE id = $2 AND (rights IS NULL OR rights::text = '{}')`

	res, err := r.db.ExecContext(ctx, query, rightsJSON, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("user already has rights, update denied")
	}

	return nil
}

func (r *userRepository) EditRightsToUser(ctx context.Context, id uint, rights map[string][]string) error {
	rightsJSON, err := json.Marshal(rights)
	if err != nil {
		return err
	}

	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM users_rights WHERE id = $1)`
	err = r.db.QueryRowContext(ctx, checkQuery, id).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("user rights with given id not found")
	}

	query := `UPDATE users_rights SET rights = $1 WHERE id = $2`
	_, err = r.db.ExecContext(ctx, query, rightsJSON, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) RevokeRightsFromUser(ctx context.Context, id uint, rights map[string][]string) error {
	var currentRightsJSON []byte
	query := `SELECT rights FROM users_rights WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&currentRightsJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user rights with given id not found")
		}
		return err
	}

	currentRights := make(map[string][]string)
	if len(currentRightsJSON) > 0 {
		if err := json.Unmarshal(currentRightsJSON, &currentRights); err != nil {
			return err
		}
	}

	for section, perms := range rights {
		existingPerms, exists := currentRights[section]
		if !exists {
			return fmt.Errorf("section %s does not exist", section) // Ошибка, если секции нет
		}

		if len(perms) == 0 {
			delete(currentRights, section)
			continue
		}

		newPerms := []string{}
		for _, existingPerm := range existingPerms {
			found := false
			for _, permToRemove := range perms {
				if existingPerm == permToRemove {
					found = true
					break
				}
			}
			if !found {
				newPerms = append(newPerms, existingPerm)
			}
		}

		if len(newPerms) == 0 {
			delete(currentRights, section)
		} else {
			currentRights[section] = newPerms
		}
	}

	updatedRightsJSON, err := json.Marshal(currentRights)
	if err != nil {
		return err
	}

	updateQuery := `UPDATE users_rights SET rights = $1 WHERE id = $2`
	_, err = r.db.ExecContext(ctx, updateQuery, updatedRightsJSON, id)
	return err
}

func (r *userRepository) GetUser(ctx context.Context, id uint) (string, string, map[string][]string, error) {
	query := `SELECT name, context, rights FROM users_rights WHERE id = $1`

	var name, context string
	var rightsJSON []byte
	err := r.db.QueryRowContext(ctx, query, id).Scan(&name, &context, &rightsJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", nil, errors.New("user not found")
		}
		return "", "", nil, err
	}

	rights := make(map[string][]string)
	if len(rightsJSON) > 0 {
		if err := json.Unmarshal(rightsJSON, &rights); err != nil {
			return "", "", nil, err
		}
	}

	return name, context, rights, nil
}
