package transport

import (
	"context"
	"encoding/json"
	"net/http"

	//"gitlab.fortebank.com/forte-market/apps/user-profile-api/src/authenticate/service"
	"github.com/rafaceo/go-test-auth/cmd/service"
	//"gitlab.fortebank.com/forte-market/apps/common-libs/errors"
	"github.com/rafaceo/go-test-auth/cmd/common-libs/errors"
)

type LoginRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success      bool   `json:"success"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Message      string `json:"message"`
}

func DecodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.InvalidCharacter.SetDevMessage(err.Error())
	}
	return req, nil
}

func EncodeLoginResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

func MakeLoginEndpoint(s service.AuthService) func(ctx context.Context, request interface{}) (interface{}, error) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LoginRequest)

		accessToken, refreshToken, err := s.Login(ctx, req.Phone, req.Password)
		if err != nil {
			return LoginResponse{Success: false, Message: err.Error()}, nil
		}

		return LoginResponse{
			Success:      true,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			Message:      "Login successful",
		}, nil
	}
}
