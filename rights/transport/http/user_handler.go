package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	kittransport "github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/rafaceo/go-test-auth/cmd/errors_auth/encoders"
	"github.com/rafaceo/go-test-auth/common-libs/httphandlers"
	"github.com/rafaceo/go-test-auth/rights/service"
	"net/http"
	"strconv"
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
	}
}

type CreateRequest struct {
	Name string `json:"name"`
}

type CreateResponse struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type EditUserRequest struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Context string `json:"context"`
}

type EditUserResponse struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type GrantRightsRequest struct {
	Rights map[string][]string `json:"rights"`
	ID     uint                `json:"id"`
}

type GrantRightsResponse struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type EditRightsRequest struct {
	Rights map[string][]string `json:"rights"`
	ID     uint                `json:"id"`
}

type EditRightsResponse struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type RevokeRightsRequest struct {
	Rights map[string][]string `json:"rights"`
	ID     uint                `json:"id"`
}

type RevokeRightsResponse struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type GetUserRequest struct {
	ID uint `json:"id"`
}

type GetUserResponse struct {
	Name    string              `json:"name,omitempty"`
	Context string              `json:"context,omitempty"`
	Rights  map[string][]string `json:"rights,omitempty"`
	Error   string              `json:"error,omitempty"`
}

func MakeCreateEndpoint(svc service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(CreateRequest)
		if !ok {
			return CreateResponse{Error: "invalid request"}, nil
		}

		err := svc.CreateUser(ctx, req.Name)
		if err != nil {
			return CreateResponse{Error: err.Error()}, nil
		}
		return CreateResponse{Message: "GoodCreated"}, nil
	}
}

func MakeEditUserEndpoint(svc service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(EditUserRequest)
		if !ok {
			return EditUserResponse{Error: "invalid request"}, nil
		}

		// Теперь передаем ID, name и context
		err := svc.EditUser(ctx, req.ID, req.Name, req.Context)
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

		name, context, rights, err := svc.GetUser(ctx, req.ID)
		if err != nil {
			return GetUserResponse{Error: err.Error()}, nil
		}

		return GetUserResponse{
			Name:    name,
			Context: context,
			Rights:  rights,
		}, nil
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

	// Используем mux.Vars, чтобы извлечь параметры из пути
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %v", err)
	}
	req.ID = uint(id)

	// Извлекаем данные из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func DecodeGrantRightsToUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req GrantRightsRequest

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %v", err)
	}
	req.ID = uint(id)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func DecodeEditRightsToUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req EditRightsRequest

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %v", err)
	}
	req.ID = uint(id)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func DecodeRevokeRightsFromUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req RevokeRightsRequest

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %v", err)
	}
	req.ID = uint(id)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func DecodeGetUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req GetUserRequest

	// Получаем ID из URL-параметров
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %v", err)
	}
	req.ID = uint(id)

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
