package transport

import (
	"github.com/rafaceo/go-test-auth/user_contexts/domain"
)

// структура ответа AddRole
type AddUserContextResponse struct {
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

// структура ответа EditRole
type EditUserContextResponse struct {
	Error string `json:"error,omitempty"`
}

// структура ответа GetRoles
type GetUserContextResponse struct {
	UserContext domain.UserContext `json:"user_context,omitempty"`
	Error       string             `json:"error,omitempty"`
}

// структура ответа DeleteRole
type DeleteUserContextResponse struct {
	Error string `json:"error,omitempty"`
}

// структура ответа DeleteRole
type DeleteAllUserContextsResponse struct {
	Error string `json:"error,omitempty"`
}
