package domain

import "time"

type Right struct {
	ID        string    `json:"id"`
	Module    string    `json:"module"`
	Action    []string  `json:"action"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
