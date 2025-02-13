package https

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	kittransport "github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/rafaceo/go-test-auth/cmd/errors_auth/encoders"
	"github.com/rafaceo/go-test-auth/common-libs/httphandlers"
	"log"
	"net/http"

	e "github.com/rafaceo/go-test-auth/cmd/errors_auth"
	"github.com/rafaceo/go-test-auth/cmd/service"
)

func GetAuthHandlers(serv service.AuthService, logger kitlog.Logger) []*httphandlers.HTTPHandler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(kittransport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encoders.EncodeErrorAUTH),
	}

	registerHandler := kithttp.NewServer(
		MakeRegisterEndpoint(serv),
		DecodeRegisterRequest,
		EncodeResponse,
		opts...,
	)

	loginHandler := kithttp.NewServer(
		MakeLoginEndpoint(serv),
		DecodeLoginRequest,
		EncodeResponse,
		opts...,
	)

	logoutHandler := kithttp.NewServer(
		MakeLogoutEndpoint(serv),
		DecodeLogoutRequest,
		EncodeResponse,
		opts...,
	)

	refreshHandler := kithttp.NewServer(
		MakeRefreshTokenEndpoint(serv),
		DecodeRefreshTokenRequest,
		EncodeResponse,
		opts...,
	)

	return []*httphandlers.HTTPHandler{
		{
			Path:    "/api/v4/users",
			Handler: registerHandler,
			Methods: []string{"POST"},
		},
		{
			Path:    "/api/v4/auth/login",
			Handler: loginHandler,
			Methods: []string{"POST"},
		},
		{
			Path:    "/api/v4/auth/logout",
			Handler: logoutHandler,
			Methods: []string{"POST"},
		},
		{
			Path:    "/api/v4/auth/refresh",
			Handler: refreshHandler,
			Methods: []string{"POST"},
		},
	}
}

type RegisterRequest struct {
	Phone     string `json:"phone"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type RegisterResponse struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type LoginRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Error        string `json:"error,omitempty"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type LogoutResponse struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Error        string `json:"error,omitempty"`
}

// MakeLoginEndpoint создаёт эндпоинт для логина
func MakeLoginEndpoint(svc service.AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LoginRequest)
		accessToken, refreshToken, err := svc.Login(ctx, req.Phone, req.Password)
		if err != nil {
			if errors.Is(err, e.TooManyRequestError) {
				return nil, err
			}
			return nil, e.UnauthorizedError
		}
		return LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
	}
}

func MakeLogoutEndpoint(svc service.AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LogoutRequest)
		err := svc.Logout(ctx, req.RefreshToken)
		if err != nil {
			return nil, e.Forbidden
		}
		return LogoutResponse{Message: "Logout successful"}, nil
	}
}
func MakeRegisterEndpoint(svc service.AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(RegisterRequest)
		if !ok {
			return RegisterResponse{Error: "invalid request"}, nil
		}

		message, err := svc.Register(ctx, req.Phone, req.Email, req.Password, req.FirstName, req.LastName)
		if err != nil {
			return RegisterResponse{Error: err.Error()}, nil
		}
		return RegisterResponse{Message: message}, nil
	}
}

func MakeRefreshTokenEndpoint(svc service.AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		log.Printf("Received request: %#v", request)
		req, ok := request.(RefreshRequest)
		if !ok {
			return RefreshResponse{Error: "invalid request"}, nil
		}
		log.Printf("Parsed RefreshRequest: %+v", req)

		accessToken, refreshToken, err := svc.RefreshToken(ctx, req.RefreshToken)
		if err != nil {
			return nil, e.Forbidden
		}
		return RefreshResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
	}
}

// DecodeLoginRequest декодирует JSON-запрос
func DecodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Phone == "" || req.Password == "" {
		return nil, e.BadRequestError
	}

	return req, nil
}

func DecodeLogoutRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.RefreshToken == "" {
		return nil, errors.New("invalid request: missing refresh_token")
	}
	return req, nil
}

func DecodeRegisterRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func DecodeRefreshTokenRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.RefreshToken == "" {
		return nil, errors.New("invalid request: missing refresh_token")
	}
	return req, nil
}

// EncodeResponse кодирует JSON-ответ
func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	prettyJSON, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return err
	}

	_, err = w.Write(prettyJSON)
	return err
}
