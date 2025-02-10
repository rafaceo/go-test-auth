package domain

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	conf "github.com/rafaceo/go-test-auth/config"
	"time"
)

type Claims struct {
	UserID       uuid.UUID `json:"id"`
	Phone        string    `json:"phone"`
	Roles        []string  `json:"roles"`
	Entitlements []string  `json:"entitlements"`
	jwt.RegisteredClaims
}

func getAccessTokenExpiration() time.Duration {
	return time.Duration(conf.AllConfigs.Env.AccessTokenExpMin) * time.Minute
}

func getRefreshTokenExpiration() time.Duration {
	return time.Duration(conf.AllConfigs.Env.RefreshTokenExpMin) * time.Minute
}

// Генерация access_token
func GenerateAccessToken(userID uuid.UUID, phone string, roles, entitlements []string, jwtSecret string) (string, error) {
	claims := Claims{
		UserID:       userID,
		Phone:        phone,
		Roles:        roles,
		Entitlements: entitlements,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(getAccessTokenExpiration())),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

// Генерация refresh_token (простая строка UUID)
func GenerateRefreshToken() string {
	return uuid.New().String()
}
