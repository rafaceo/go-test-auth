package domain

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"os"
	"time"
)

const (
	AccessTokenExpiration  = time.Hour * 1        // Живет 1 час
	RefreshTokenExpiration = (time.Hour * 24) * 7 // Живет 7 дней
)

type Claims struct {
	UserID       uuid.UUID `json:"id"`
	Phone        string    `json:"phone"`
	Roles        []string  `json:"roles"`
	Entitlements []string  `json:"entitlements"`
	jwt.RegisteredClaims
}

// Генерация access_token
func GenerateAccessToken(userID uuid.UUID, phone string, roles, entitlements []string) (string, error) {
	claims := Claims{
		UserID:       userID,
		Phone:        phone,
		Roles:        roles,
		Entitlements: entitlements,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// Генерация refresh_token (простая строка UUID)
func GenerateRefreshToken() string {
	return uuid.New().String()
}
