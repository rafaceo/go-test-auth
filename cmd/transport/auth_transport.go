package https

//
//import (
//	"context"
//	"encoding/json"
//	"errors"
//	"github.com/go-kit/kit/endpoint"
//	kitlog "github.com/go-kit/kit/log"
//	kittransport "github.com/go-kit/kit/transport"
//	httptransport "github.com/go-kit/kit/transport/http"
//	kithttp "github.com/go-kit/kit/transport/http"
//	"github.com/rafaceo/go-test-auth/common-libs/encoders"
//	"github.com/rafaceo/go-test-auth/common-libs/httphandlers"
//	"net/http"
//
//	"github.com/rafaceo/go-test-auth/cmd/service"
//
//
//)
//
//func GetAuthHandlers(serv service.AuthService, logger kitlog.Logger) []*httphandlers.HTTPHandler {
//	opts := []kithttp.ServerOption{
//		kithttp.ServerErrorHandler(kittransport.NewLogErrorHandler(logger)),
//		kithttp.ServerErrorEncoder(encoders.EncodeErrorJSON),
//	}
//
//	loginHandler := kithttp.NewServer(
//		MakeLoginEndpoint(serv),
//		DecodeLoginRequest,
//		EncodeResponse,
//		opts...,
//	)
//
//	logoutHandler := kithttp.NewServer(
//		MakeLogoutEndpoint(serv),
//		DecodeLogoutRequest,
//		EncodeResponse,
//		opts...,
//	)
//
//	refreshHandler := kithttp.NewServer(
//		MakeRefreshTokenEndpoint(serv),
//		DecodeRefreshTokenRequest,
//		EncodeResponse,
//		opts...,
//	)
//
//	return []*httphandlers.HTTPHandler{
//		{
//			Path:    "/api/v4/auth/login",
//			Handler: loginHandler,
//			Methods: []string{"POST"},
//		},
//		{
//			Path:    "/api/v4/auth/logout",
//			Handler: logoutHandler,
//			Methods: []string{"POST"},
//		},
//		{
//			Path:    "/api/v4/auth/refresh",
//			Handler: refreshHandler,
//			Methods: []string{"POST"},
//		},
//	}
//}
//
//type LoginRequest struct {
//	Phone    string `json:"phone"`
//	Password string `json:"password"`
//}
//
//type LogoutRequest struct {
//	RefreshToken string `json:"refresh_token"`
//}
//
//type LogoutResponse struct {
//	Message string `json:"message,omitempty"`
//	Error   string `json:"error,omitempty"`
//}
//
//type RefreshRequest struct {
//	RefreshToken string `json:"refresh_token"`
//}
//
//type RefreshResponse struct {
//	AccessToken  string `json:"access_token,omitempty"`
//	RefreshToken string `json:"refresh_token,omitempty"`
//	Error        string `json:"error,omitempty"`
//}
//
//// LoginResponse структура ответа
//type LoginResponse struct {
//	AccessToken  string `json:"access_token,omitempty"`
//	RefreshToken string `json:"refresh_token,omitempty"`
//	Error        string `json:"error,omitempty"`
//}
//
//// MakeLoginEndpoint создаёт эндпоинт для логина
//func MakeLoginEndpoint(svc service.AuthService) endpoint.Endpoint {
//	return func(ctx context.Context, request interface{}) (interface{}, error) {
//		req := request.(LoginRequest)
//		accessToken, refreshToken, err := svc.Login(ctx, req.Phone, req.Password)
//		if err != nil {
//			return LoginResponse{Error: err.Error()}, nil
//		}
//		return LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
//	}
//}
//
//func MakeLogoutEndpoint(svc service.AuthService) endpoint.Endpoint {
//	return func(ctx context.Context, request interface{}) (interface{}, error) {
//		req := request.(LogoutRequest)
//		err := svc.Logout(ctx, req.RefreshToken)
//		if err != nil {
//			return LogoutResponse{Error: err.Error()}, nil
//		}
//		return LogoutResponse{Message: "Logout successful"}, nil
//	}
//}
//
//func MakeRefreshTokenEndpoint(svc service.AuthService) endpoint.Endpoint {
//	return func(ctx context.Context, request interface{}) (interface{}, error) {
//		req := request.(RefreshRequest)
//		accessToken, refreshToken, err := svc.RefreshToken(ctx, req.RefreshToken)
//		if err != nil {
//			return RefreshResponse{Error: err.Error()}, nil
//		}
//		return RefreshResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
//	}
//}
//
//// DecodeLoginRequest декодирует JSON-запрос
//func DecodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
//	var req LoginRequest
//	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
//		return nil, err
//	}
//	return req, nil
//}
//
//func DecodeLogoutRequest(_ context.Context, r *http.Request) (interface{}, error) {
//	var req LogoutRequest
//	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.RefreshToken == "" {
//		return nil, errors.New("invalid request: missing refresh_token")
//	}
//	return req, nil
//}
//
//func DecodeRefreshTokenRequest(_ context.Context, r *http.Request) (interface{}, error) {
//	var req RefreshRequest
//	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.RefreshToken == "" {
//		return nil, errors.New("invalid request")
//	}
//	return req, nil
//}
//
//// EncodeResponse кодирует JSON-ответ
//func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
//	return json.NewEncoder(w).Encode(response)
//}
//
//// NewLoginHandler создаёт HTTP-обработчик для логина
//func NewLoginHandler(authService service.AuthService) http.Handler {
//	return httptransport.NewServer(
//		MakeLoginEndpoint(authService),
//		DecodeLoginRequest,
//		EncodeResponse,
//	)
//}
//
//func NewLogoutHandler(authService service.AuthService) http.Handler {
//	return httptransport.NewServer(
//		MakeLogoutEndpoint(authService),
//		DecodeLogoutRequest,
//		EncodeResponse,
//	)
//}
//
//func NewRefreshHandler(authService service.AuthService) http.Handler {
//	return httptransport.NewServer(
//		MakeRefreshTokenEndpoint(authService),
//		DecodeRefreshTokenRequest,
//		EncodeResponse,
//	)
//}
