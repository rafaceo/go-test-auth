package domain

import (
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID
	Phone    string
	Password string
}
