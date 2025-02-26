package transport

import (
	"github.com/rafaceo/go-test-auth/roles/domain"
)

type AddRoleResponse struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type EditRoleResponse struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type GetRolesResponse struct {
	Roles []domain.Role `json:"roles"`
	Error string        `json:"error,omitempty"`
}

type GetRoleRightsResponse struct {
	Rights map[string][]string `json:"rights"`
	Error  string              `json:"error,omitempty"`
}

type DeleteRoleResponse struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type AssignRoleToUserResponse struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}
