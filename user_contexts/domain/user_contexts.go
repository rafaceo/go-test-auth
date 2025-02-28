package domain

import "github.com/google/uuid"

type UserContext struct {
	UserID     uuid.UUID `json:"user_id"`
	MerchantID string    `json:"merchant_id"`
	Global     bool      `json:"global"`
}
