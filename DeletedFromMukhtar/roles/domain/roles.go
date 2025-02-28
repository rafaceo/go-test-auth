package domain

import (
	"github.com/google/uuid"
	"time"
)

type Role struct {
	RoleID    uuid.UUID `json:"role_id"`
	RoleName  string    `json:"role_name"`
	Rights    []string  `json:"rights"` // JSONB будет храниться как массив строк
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
