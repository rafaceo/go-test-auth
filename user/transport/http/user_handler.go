package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	kittransport "github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rafaceo/go-test-auth/cmd/errors_auth/encoders"
	"github.com/rafaceo/go-test-auth/common-libs/httphandlers"
	"github.com/rafaceo/go-test-auth/user/service"
	"net/http"
)

func GetUserHandler(serv service.UserService, logger kitlog.Logger) []*httphandlers.HTTPHandler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(kittransport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encoders.EncodeErrorAUTH),
	}

	createUser := kithttp.NewServer(
		MakeCreateEndpoint(serv),
		DecodeCreateRequest,
		EncodeResponse,
		opts...,
	)

	editUser := kithttp.NewServer(
		MakeEditUserEndpoint(serv),
		DecodeEditUserRequest,
		EncodeResponse,
		opts...,
	)

	grantRights := kithttp.NewServer(
		MakeGrantRightsToUserEndpoint(serv),
		DecodeGrantRightsToUserRequest,
		EncodeResponse,
		opts...,
	)

	editRights := kithttp.NewServer(
		MakeEditRightsToUserEndpoint(serv),
		DecodeEditRightsToUserRequest,
		EncodeResponse,
		opts...,
	)

	revokeRights := kithttp.NewServer(
		MakeRevokeRightsFromUserEndpoint(serv),
		DecodeRevokeRightsFromUserRequest,
		EncodeResponse,
		opts...,
	)

	getUser := kithttp.NewServer(
		MakeGetUserEndpoint(serv),
		DecodeGetUserRequest,
		EncodeResponse,
		opts...,
	)

	getUserRights := kithttp.NewServer(
		MakeGetUserRightsEndpoint(serv),
		DecodeGetUserRightsRequest,
		EncodeResponse,
		opts...,
	)

	return []*httphandlers.HTTPHandler{
		{
			Path:    "/api/v4/userss",
			Handler: createUser,
			Methods: []string{"POST"},
		},
		{
			Path:    "/api/v4/users/{id}",
			Handler: editUser,
			Methods: []string{"PUT"},
		},
		{
			Path:    "/api/v4/users/{id}/rights",
			Handler: grantRights,
			Methods: []string{"POST"},
		},
		{
			Path:    "/api/v4/users/{id}/rights",
			Handler: editRights,
			Methods: []string{"PUT"},
		},
		{
			Path:    "/api/v4/users/{id}/rights",
			Handler: revokeRights,
			Methods: []string{"DELETE"},
		},
		{
			Path:    "/api/v4/userss/{id}",
			Handler: getUser,
			Methods: []string{"GET"},
		},
		{
			Path:    "/api/v4/users/{id}/rights",
			Handler: getUserRights,
			Methods: []string{"GET"},
		},
	}
}

type CreateRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type CreateResponse struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type EditUserRequest struct {
	ID       uuid.UUID
	Phone    string
	Password string
}

type EditUserResponse struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type GrantRightsRequest struct {
	Rights map[string][]string `json:"rights"`
	ID     uuid.UUID           `json:"id"`
}

type GrantRightsResponse struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type EditRightsRequest struct {
	Rights map[string][]string `json:"rights"`
	ID     uuid.UUID           `json:"id"`
}

type EditRightsResponse struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type RevokeRightsRequest struct {
	Rights map[string][]string `json:"rights"`
	ID     uuid.UUID           `json:"id"`
}

type RevokeRightsResponse struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type GetUserRequest struct {
	ID uuid.UUID `json:"id"`
}

type GetUserResponse struct {
	Phone     string `json:"phone,omitempty"`
	Password  string `json:"password,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	Error     string `json:"error,omitempty"`
}

type GetUserRightsRequest struct {
	ID uuid.UUID `json:"id"`
}

type GetUserRightsResponse struct {
	Rights map[string][]string `json:"rights"`
	Error  string              `json:"error,omitempty"`
}

func MakeCreateEndpoint(svc service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(CreateRequest)
		if !ok {
			return CreateResponse{Error: "invalid request"}, nil
		}

		err := svc.CreateUser(ctx, req.Phone, req.Password)
		if err != nil {
			return CreateResponse{Error: err.Error()}, nil
		}
		return CreateResponse{Message: "Пользователь успешно зарегистрирован"}, nil
	}
}

func MakeEditUserEndpoint(svc service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(EditUserRequest)
		if !ok {
			return EditUserResponse{Error: "invalid request"}, nil
		}

		// Теперь передаем ID, name и context
		err := svc.EditUser(ctx, req.ID, req.Phone, req.Password)
		if err != nil {
			return EditUserResponse{Error: err.Error()}, nil
		}
		return EditUserResponse{Message: "User updated successfully"}, nil
	}
}

func MakeGrantRightsToUserEndpoint(svc service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(GrantRightsRequest)
		if !ok {
			return GrantRightsResponse{Error: "invalid request"}, nil
		}

		err := svc.GrantRightsToUser(ctx, req.ID, req.Rights)
		if err != nil {
			return GrantRightsResponse{Error: err.Error()}, nil
		}
		return GrantRightsResponse{Message: "Rights granted successfully"}, nil
	}
}

func MakeEditRightsToUserEndpoint(svc service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(EditRightsRequest)
		if !ok {
			return EditRightsResponse{Error: "invalid request"}, nil
		}

		err := svc.EditRightsToUser(ctx, req.ID, req.Rights)
		if err != nil {
			return EditRightsResponse{Error: err.Error()}, nil
		}
		return EditRightsResponse{Message: "Rights changed successfully"}, nil
	}
}

func MakeRevokeRightsFromUserEndpoint(svc service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(RevokeRightsRequest)
		if !ok {
			return RevokeRightsResponse{Error: "invalid request"}, nil
		}

		err := svc.RevokeRightsFromUser(ctx, req.ID, req.Rights)
		if err != nil {
			return RevokeRightsResponse{Error: err.Error()}, nil
		}
		return RevokeRightsResponse{Message: "Rights revoked successfully"}, nil
	}
}

func MakeGetUserEndpoint(svc service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(GetUserRequest)
		if !ok {
			return GetUserResponse{Error: "invalid request"}, nil
		}

		phone, passwordHash, createdAt, updatedAt, err := svc.GetUser(ctx, req.ID)
		if err != nil {
			return GetUserResponse{Error: err.Error()}, nil
		}

		return GetUserResponse{
			Phone:     phone,
			Password:  passwordHash,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}, nil
	}
}

func MakeGetUserRightsEndpoint(svc service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(GetUserRightsRequest)
		if !ok {
			return GetUserRightsResponse{Error: "invalid request"}, nil
		}

		rights, err := svc.GetUserRights(ctx, req.ID)
		if err != nil {
			return GetUserRightsResponse{Error: err.Error()}, nil
		}

		return GetUserRightsResponse{Rights: rights}, nil
	}
}

func DecodeCreateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func DecodeEditUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req EditUserRequest

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID: %v", err)
	}
	req.ID = id

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("invalid request body: %v", err)
	}

	return req, nil
}

func DecodeGrantRightsToUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req GrantRightsRequest

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID: %v", err)
	}
	req.ID = id

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("invalid request body: %v", err)
	}

	return req, nil
}

func DecodeEditRightsToUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req EditRightsRequest

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID: %v", err)
	}
	req.ID = id

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("invalid request body: %v", err)
	}

	return req, nil
}

func DecodeRevokeRightsFromUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req RevokeRightsRequest

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID: %v", err)
	}
	req.ID = id

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("invalid request body: %v", err)
	}

	return req, nil
}

func DecodeGetUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req GetUserRequest

	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		return nil, fmt.Errorf("missing user ID in URL")
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID: %v", err)
	}
	req.ID = id

	return req, nil
}

func DecodeGetUserRightsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req GetUserRightsRequest

	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		return nil, fmt.Errorf("missing user ID in URL")
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID: %v", err)
	}
	req.ID = id

	return req, nil
}

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
