package domain

import (
	"database/sql"
	"github.com/google/uuid"
)

type UserAuth struct {
	ID       uuid.UUID `db:"id"`
	Phone    string    `db:"phone"`
	Password string    `db:"password"`

	CreatedAt sql.NullTime `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}
