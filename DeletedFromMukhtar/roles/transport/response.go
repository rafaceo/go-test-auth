package transport

import (
	"github.com/google/uuid"
	"gitlab.fortebank.com/forte-market/apps/user-profile-api/src/roles/domain"
)

// структура ответа AddRole
type AddRoleResponse struct {
	ID    uuid.UUID `json:"id,omitempty"`
	Error string    `json:"error,omitempty"`
}

// структура ответа EditRole
type EditRoleResponse struct {
	Error string `json:"error,omitempty"`
}

// структура ответа GetRoles
type GetRolesResponse struct {
	Role  domain.Role   `json:"Role,omitempty"`
	Roles []domain.Role `json:"Roles,omitempty"` // Массив объектов
	Error string        `json:"error,omitempty"`
}

// структура ответа DeleteRole
type DeleteRoleResponse struct {
	Error string `json:"error,omitempty"`
}
