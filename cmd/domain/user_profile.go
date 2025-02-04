package domain

import (
	"database/sql"
	"github.com/google/uuid"
)

type UserProfile struct {
	ID           uuid.UUID `db:"id"`
	IIN          string    `db:"iin"`
	Name         string    `db:"name"`
	Email        string    `db:"email"`
	Phone        string    `db:"phone"`
	Subscribe    bool      `db:"subscribe"`
	SourceAuth   string    `db:"source_auth"`
	Roles        []string  `db:"roles"`
	Entitlements []string  `db:"entitlements"`
	Password     string    `db:"password"`

	RefreshToken string       `db:"refresh_token"`
	CreatedAt    sql.NullTime `db:"created_at"`
	UpdatedAt    sql.NullTime `db:"updated_at"`
}
