package service

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
	dom "github.com/rafaceo/go-test-auth/cmd/domain"
	"net/http"
	"time"
)

func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req dom.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := jwt.Parse(req.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return dom.SecretKey, nil
	})

	if err != nil || !token.Valid {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	exp, ok := claims["exp"].(float64)
	if !ok || time.Now().Unix() > int64(exp) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	expirationTime := time.Now().Add(15 * time.Minute).Unix()
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": expirationTime,
	})
	accessToken, err := newToken.SignedString(dom.SecretKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dom.RefreshResponse{AccessToken: accessToken})
}
