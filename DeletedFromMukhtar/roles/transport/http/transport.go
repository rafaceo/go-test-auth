package http

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	kittransport "github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	kithttp "github.com/go-kit/kit/transport/http"
	"gitlab.fortebank.com/forte-market/apps/common-libs/encoders"
	"gitlab.fortebank.com/forte-market/apps/common-libs/httphandlers"
	service "gitlab.fortebank.com/forte-market/apps/user-profile-api/src/roles/service"
	"gitlab.fortebank.com/forte-market/apps/user-profile-api/src/roles/transport"
	"net/http"
)

func GetRoleHandlers(serv service.RoleService, logger kitlog.Logger) []*httphandlers.HTTPHandler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(kittransport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encoders.EncodeErrorJSON),
	}

	addRolesHandler := kithttp.NewServer(
		MakeAddRoleEndpoint(serv),
		DecodeAddRoleRequest,
		EncodeResponse,
		opts...,
	)

	editRolesHandler := kithttp.NewServer(
		MakeEditRoleEndpoint(serv),
		DecodeEditRoleRequest,
		EncodeResponse,
		opts...,
	)

	getRolesByIdHandler := kithttp.NewServer(
		MakeGetRolesEndpoint(serv),
		DecodeGetRoleRequest,
		EncodeResponse,
		opts...,
	)

	getAllRolesHandler := kithttp.NewServer(
		MakeGetRolesEndpoint(serv),
		DecodeGetRoleRequest,
		EncodeResponse,
		opts...,
	)

	deleteRolesHandler := kithttp.NewServer(
		MakeDeleteRoleEndpoint(serv),
		DecodeDeleteRoleRequest,
		EncodeResponse,
		opts...,
	)
	return []*httphandlers.HTTPHandler{
		{
			Path:    "/api/v4/roles",
			Handler: addRolesHandler,
			Methods: []string{"POST"},
		},
		{
			Path:    "/api/v4/roles/{Roles_id}",
			Handler: editRolesHandler,
			Methods: []string{"PUT"},
		},
		{
			Path:    "/api/v4/roles/{Roles_id}",
			Handler: getRolesByIdHandler,
			Methods: []string{"GET"},
		},
		{
			Path:    "/api/v4/roles",
			Handler: getAllRolesHandler,
			Methods: []string{"GET"},
		},
		{
			Path:    "/api/v4/roles/{Roles_id}",
			Handler: deleteRolesHandler,
			Methods: []string{"DELETE"},
		},
	}
}

// MakeAddRoleEndpoint создаёт эндпоинт для добавления прав
func MakeAddRoleEndpoint(svc service.RoleService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(transport.AddRoleRequest)
		id, err := svc.AddRole(ctx, req.Role)
		if err != nil {
			return transport.AddRoleResponse{Error: err.Error()}, nil
		}
		return transport.AddRoleResponse{id, ""}, nil
	}
}

// MakeEditRoleEndpoint создаёт эндпоинт для изменения прав
func MakeEditRoleEndpoint(svc service.RoleService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(transport.EditRoleRequest)
		err := svc.UpdateRole(ctx, req.Role)
		if err != nil {
			return transport.EditRoleResponse{Error: err.Error()}, nil
		}
		return transport.EditRoleResponse{}, nil
	}
}

// MakeGetRolesEndpoint создаёт эндпоинт для вывода прав
func MakeGetRolesEndpoint(svc service.RoleService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_, ok := request.(transport.GetRoleRequest)
		if !ok {
			return nil, errors.New("invalid request type")
		}

		// Получаем все роли, если ID не передан
		roles, err := svc.GetAllRoles(ctx)
		if err != nil {
			return transport.GetRolesResponse{Error: err.Error()}, nil
		}
		return transport.GetRolesResponse{Roles: roles}, nil
	}
}

// MakeDeleteRoleEndpoint создаёт эндпоинт для удаления прав
func MakeDeleteRoleEndpoint(svc service.RoleService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(transport.DeleteRoleRequest)
		err := svc.DeleteRole(ctx, req.ID)
		if err != nil {
			return transport.DeleteRoleResponse{Error: err.Error()}, nil
		}
		return transport.DeleteRoleResponse{}, nil
	}
}

// DecodeAddRoleRequest декодирует JSON-запрос
func DecodeAddRoleRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req transport.AddRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

// DecodeEditRoleRequest декодирует JSON-запрос
func DecodeEditRoleRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req transport.EditRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

// DecodeGetRoleRequest декодирует JSON-запрос
func DecodeGetRoleRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req transport.EditRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

// DecodeDeleteRoleRequest декодирует JSON-запрос
func DecodeDeleteRoleRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req transport.DeleteRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

// EncodeResponse кодирует JSON-ответ
func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

// NewLoginHandler создаёт HTTP-обработчик для логина
func NewAddRolesHandler(RolesService service.RoleService) http.Handler {
	return httptransport.NewServer(
		MakeAddRoleEndpoint(RolesService),
		DecodeAddRoleRequest,
		EncodeResponse,
	)
}
