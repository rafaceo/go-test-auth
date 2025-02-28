package transport

import (
	"github.com/google/uuid"
	"gitlab.fortebank.com/forte-market/apps/user-profile-api/src/roles/domain"
)

type AddRoleRequest struct {
	Role domain.Role `json:"Role,omitempty"`
}

type EditRoleRequest struct {
	Role domain.Role `json:"Role,omitempty"`
}

type GetRoleRequest struct {
	ID uuid.UUID `json:"id,omitempty"`
}

type DeleteRoleRequest struct {
	ID uuid.UUID `json:"id,omitempty"`
}
