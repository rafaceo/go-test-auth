package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	kittransport "github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/rafaceo/go-test-auth/cmd/errors_auth/encoders"
	"github.com/rafaceo/go-test-auth/common-libs/httphandlers"
	"github.com/rafaceo/go-test-auth/roles/domain"
	service "github.com/rafaceo/go-test-auth/roles/service"
	"strconv"

	"github.com/rafaceo/go-test-auth/roles/transport"
	"net/http"
)

func GetRoleHandlers(serv service.RoleService, logger kitlog.Logger) []*httphandlers.HTTPHandler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(kittransport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encoders.EncodeErrorAUTH),
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

	getAllRolesHandler := kithttp.NewServer(
		MakeGetRolesEndpoint(serv),
		DecodeGetRolesRequest,
		EncodeResponse,
		opts...,
	)

	getRolesRightsById := kithttp.NewServer(
		MakeGetRoleRightsEndpoint(serv),
		DecodeGetRoleRightsRequest,
		EncodeResponse,
		opts...,
	)

	deleteRolesHandler := kithttp.NewServer(
		MakeDeleteRoleEndpoint(serv),
		DecodeDeleteRoleRequest,
		EncodeResponse,
		opts...,
	)

	assignRoleToUserHandler := kithttp.NewServer(
		MakeAssignRoleToUserEndpoint(serv),
		DecodeAssignRoleToUserRequest,
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
			Path:    "/api/v4/roles/{role_id}",
			Handler: editRolesHandler,
			Methods: []string{"PUT"},
		},
		{
			Path:    "/api/v4/roles/{role_id}/rights",
			Handler: getRolesRightsById,
			Methods: []string{"GET"},
		},
		{
			Path:    "/api/v4/roles",
			Handler: getAllRolesHandler,
			Methods: []string{"GET"},
		},
		{
			Path:    "/api/v4/roles/{role_id}",
			Handler: deleteRolesHandler,
			Methods: []string{"DELETE"},
		},
		{
			Path:    "/api/v4/users/{role_id}/role-preset",
			Handler: assignRoleToUserHandler,
			Methods: []string{"PUT"},
		},
	}
}

func MakeAddRoleEndpoint(svc service.RoleService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(transport.AddRoleRequest)

		err := svc.AddRole(ctx, req.RoleName, req.RoleNameRu, req.Notes, req.Rights)
		if err != nil {
			return transport.AddRoleResponse{Error: err.Error()}, nil
		}
		return transport.AddRoleResponse{Message: "Роль успешно добавлена"}, nil
	}
}

func MakeEditRoleEndpoint(svc service.RoleService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(transport.EditRoleRequest)

		err := svc.EditRole(ctx, req.RoleID, req.RoleName, req.RoleNameRu, req.Notes, req.Rights)
		if err != nil {
			return transport.EditRoleResponse{Error: err.Error()}, nil
		}

		return transport.EditRoleResponse{Message: "Роль успешно изменена"}, nil
	}
}

func MakeGetRolesEndpoint(svc service.RoleService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		roles, err := svc.GetRoles(ctx)
		if err != nil {
			return transport.GetRolesResponse{Error: err.Error()}, nil
		}

		var responseRoles []domain.Role
		for _, r := range roles {
			responseRoles = append(responseRoles, domain.Role{
				ID:     r.ID,
				Name:   r.Name,
				NameRu: r.NameRu,
				Notes:  r.Notes,
				Rights: r.Rights,
			})
		}

		return transport.GetRolesResponse{Roles: responseRoles}, nil
	}
}

func MakeGetRoleRightsEndpoint(svc service.RoleService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(transport.GetRoleRightsRequest)
		rights, err := svc.GetRoleRights(ctx, req.RoleID)
		if err != nil {
			return transport.GetRoleRightsResponse{Error: err.Error()}, nil
		}
		return transport.GetRoleRightsResponse{Rights: rights}, nil
	}
}

func MakeDeleteRoleEndpoint(svc service.RoleService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(transport.DeleteRoleRequest)
		err := svc.DeleteRole(ctx, req.RoleID)
		if err != nil {
			return transport.DeleteRoleResponse{Error: err.Error()}, nil
		}
		return transport.DeleteRoleResponse{Message: "Роль успешно удалена"}, nil
	}
}

func MakeAssignRoleToUserEndpoint(svc service.RoleService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(transport.AssignRoleToUserRequest)
		err := svc.AssignRoleToUser(ctx, req.UserID, req.RoleID, req.Merge)
		if err != nil {
			return transport.AssignRoleToUserResponse{Error: err.Error()}, nil
		}
		return transport.AssignRoleToUserResponse{Message: "OKAY"}, nil
	}
}

func DecodeAddRoleRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req transport.AddRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	if req.RoleName == "" {
		return nil, errors.New("role_name обязателен")
	}

	return req, nil
}

func DecodeEditRoleRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req transport.EditRoleRequest

	vars := mux.Vars(r)
	roleID, err := strconv.Atoi(vars["role_id"])
	if err != nil {
		return nil, errors.New("invalid role ID")
	}
	req.RoleID = roleID

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func DecodeGetRolesRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return transport.GetRoleRequest{}, nil
}

func DecodeGetRoleRightsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	roleID, err := strconv.Atoi(vars["role_id"])
	if err != nil {
		return nil, err
	}
	return transport.GetRoleRightsRequest{RoleID: roleID}, nil
}

func DecodeDeleteRoleRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	roleID, err := strconv.Atoi(vars["role_id"])
	if err != nil {
		return nil, err
	}

	return transport.DeleteRoleRequest{RoleID: roleID}, nil
}

func DecodeAssignRoleToUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)

	roleID, err := strconv.Atoi(vars["role_id"])
	if err != nil {
		return nil, fmt.Errorf("invalid role_id in URL: %w", err)
	}

	var req transport.AssignRoleToUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("invalid JSON body: %w", err)
	}

	req.RoleID = roleID

	return req, nil
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(response)
}

func NewAddRolesHandler(RolesService service.RoleService) http.Handler {
	return httptransport.NewServer(
		MakeAddRoleEndpoint(RolesService),
		DecodeAddRoleRequest,
		EncodeResponse,
	)
}
