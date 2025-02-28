package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	domain "github.com/rafaceo/go-test-auth/rights/domain"
	"github.com/rafaceo/go-test-auth/rights/repository"
	"log"
	"strings"
)

type PostgresRightsRepository struct {
	db *sqlx.DB
}

func NewPostgresRightsRepository(db *sqlx.DB) repository.RightsRepository {
	return &PostgresRightsRepository{db: db}
}

func (r *PostgresRightsRepository) AddRights(ctx context.Context, module string, action []string) error {
	id := uuid.New().String() // Генерация UUID

	query := `INSERT INTO rights (id, module, action, created_at, updated_at) 
	          VALUES ($1, $2, $3, now(), now())`
	_, err := r.db.ExecContext(ctx, query, id, module, pq.Array(action))
	return err
}

func (r *PostgresRightsRepository) EditRight(ctx context.Context, id string, module string, action []string) error {
	actionJSON, err := json.Marshal(action)
	if err != nil {
		return err
	}

	query := `UPDATE rights SET module = $1, action = $2, updated_at = now() WHERE id = $3`
	_, err = r.db.ExecContext(ctx, query, module, string(actionJSON), id)
	return err
}

func (r *PostgresRightsRepository) GetAllRights(ctx context.Context) ([]domain.Right, error) {
	query := `SELECT id, module, action, created_at, updated_at FROM rights`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rights []domain.Right
	for rows.Next() {
		var right domain.Right
		var actionRaw string

		err := rows.Scan(&right.ID, &right.Module, &actionRaw, &right.CreatedAt, &right.UpdatedAt)
		if err != nil {
			return nil, err
		}

		right.Action = parseAction(actionRaw)

		rights = append(rights, right)
	}

	return rights, nil
}

func (r *PostgresRightsRepository) GetRightByName(ctx context.Context, module string) (*domain.Right, error) {
	query := `SELECT id, module, action, created_at, updated_at FROM rights WHERE module = $1`
	row := r.db.QueryRowContext(ctx, query, module)

	var right domain.Right
	var actionRaw string

	err := row.Scan(&right.ID, &right.Module, &actionRaw, &right.CreatedAt, &right.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	right.Action = parseAction(actionRaw)

	return &right, nil
}

func (r *PostgresRightsRepository) GetRightById(ctx context.Context, id string) (*domain.Right, error) {
	var right domain.Right
	var actionRaw string

	query := `SELECT id, module, action, created_at, updated_at FROM rights WHERE id = $1`
	log.Println("ID: " + id)

	err := r.db.QueryRowContext(ctx, query, id).Scan(&right.ID, &right.Module, &actionRaw, &right.CreatedAt, &right.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	right.Action = parseAction(actionRaw)

	return &right, nil
}

func (r *PostgresRightsRepository) DeleteRight(ctx context.Context, id string) error {
	query := `DELETE FROM rights WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func parseAction(actionRaw string) []string {
	actionRaw = strings.ReplaceAll(actionRaw, "{", "")
	actionRaw = strings.ReplaceAll(actionRaw, "}", "")
	actionRaw = strings.ReplaceAll(actionRaw, `"`, "")

	if actionRaw == "" {
		return []string{}
	}

	return strings.Split(actionRaw, ",")
}
