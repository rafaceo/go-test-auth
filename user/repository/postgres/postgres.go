package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	repo "github.com/rafaceo/go-test-auth/user/repository"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) repo.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, phone string, passwordHash string) error {
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM "users" WHERE phone = $1)`
	err := r.db.GetContext(ctx, &exists, checkQuery, phone)
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("пользователь с таким номером уже существует")
	}

	query := `INSERT INTO "users" (id, phone, password_hash, rights, created_at) 
	          VALUES ($1, $2, $3, $4, NOW())`
	_, err = r.db.ExecContext(ctx, query, uuid.New(), phone, passwordHash, `{}`)
	return err
}

func (r *userRepository) EditUser(ctx context.Context, id uuid.UUID, phone, password string) error {
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM "users" WHERE id = $1)`
	err := r.db.QueryRowContext(ctx, checkQuery, id).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("user with given id not found")
	}

	query := `UPDATE "users" SET phone = $1, password_hash = $2 WHERE id = $3`
	_, err = r.db.ExecContext(ctx, query, phone, password, id)
	return err
}
func (r *userRepository) GrantRightsToUser(ctx context.Context, id uuid.UUID, rights map[string][]string) error {
	rightsJSON, err := json.Marshal(rights)
	if err != nil {
		return err
	}

	// Запрос обновления прав, если текущие права NULL или пустые
	query := `
		UPDATE "users" 
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

func (r *userRepository) EditRightsToUser(ctx context.Context, id uuid.UUID, rights map[string][]string) error {
	rightsJSON, err := json.Marshal(rights)
	if err != nil {
		return err
	}

	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM "users" WHERE id = $1)`
	err = r.db.QueryRowContext(ctx, checkQuery, id).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("user rights with given id not found")
	}

	query := `UPDATE "users" SET rights = $1 WHERE id = $2`
	_, err = r.db.ExecContext(ctx, query, rightsJSON, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) RevokeRightsFromUser(ctx context.Context, id uuid.UUID, rights map[string][]string) error {
	var currentRightsJSON []byte
	query := `SELECT rights FROM "users" WHERE id = $1`
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

	updateQuery := `UPDATE "users" SET rights = $1 WHERE id = $2`
	_, err = r.db.ExecContext(ctx, updateQuery, updatedRightsJSON, id)
	return err
}

func (r *userRepository) GetUser(ctx context.Context, id uuid.UUID) (string, string, string, string, error) {
	// Проверяем, что переданный UUID корректный
	if id == uuid.Nil {
		return "", "", "", "", fmt.Errorf("invalid UUID: cannot be empty")
	}

	query := `SELECT phone, password_hash, created_at, updated_at FROM "users" WHERE id = $1`

	var phone, passwordHash, createdAt, updatedAt string
	err := r.db.QueryRowContext(ctx, query, id).Scan(&phone, &passwordHash, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", "", "", fmt.Errorf("user not found")
		}
		return "", "", "", "", fmt.Errorf("database error: %w", err)
	}

	return phone, passwordHash, createdAt, updatedAt, nil
}

func (r *userRepository) GetUserRights(ctx context.Context, id uuid.UUID) (map[string][]string, error) {
	query := `SELECT rights FROM "users" WHERE id = $1`

	var rightsJSON []byte
	err := r.db.QueryRowContext(ctx, query, id).Scan(&rightsJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Декодируем JSON в map
	var rights map[string][]string
	if len(rightsJSON) > 0 {
		if err := json.Unmarshal(rightsJSON, &rights); err != nil {
			return nil, fmt.Errorf("failed to parse rights JSON: %v", err)
		}
	} else {
		rights = make(map[string][]string) // Если пусто, возвращаем пустую мапу
	}

	return rights, nil
}
