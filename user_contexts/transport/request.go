package transport

import (
	"github.com/google/uuid"
)

type AddUserContextRequest struct {
	UserID     uuid.UUID `json:"user_id,omitempty"`
	MerchantID string    `json:"merchant_id,omitempty"`
	Global     bool      `json:"global,omitempty"`
}

type EditUserContextRequest struct {
	UserID uuid.UUID `json:"user_id"`
	Global bool      `json:"global"`
}

type GetUserContextRequest struct {
	UserID uuid.UUID `json:"user_id"`
}

type DeleteUserContextRequest struct {
	UserID     uuid.UUID `json:"user_id"`
	MerchantID string    `json:"merchant_id"`
}

type DeleteAllUserContextsRequest struct {
	UserID uuid.UUID `json:"user_id"`
}
