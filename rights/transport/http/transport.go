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
	"github.com/gorilla/mux"
	"github.com/rafaceo/go-test-auth/cmd/errors_auth/encoders"
	"github.com/rafaceo/go-test-auth/common-libs/httphandlers"
	"github.com/rafaceo/go-test-auth/rights/domain"
	"github.com/rafaceo/go-test-auth/rights/service"
	"net/http"
)

func GetRightHandlers(serv service.RightsService, logger kitlog.Logger) []*httphandlers.HTTPHandler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(kittransport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encoders.EncodeErrorAUTH),
	}

	addRightHandler := kithttp.NewServer(
		MakeAddRightEndpoint(serv),
		DecodeAddRightRequest,
		EncodeResponse,
		opts...,
	)

	editRightHandler := kithttp.NewServer(
		MakeEditRightEndpoint(serv),
		DecodeEditRightRequest,
		EncodeResponse,
		opts...,
	)

	getAllRightsHandler := kithttp.NewServer(
		MakeGetAllRightsEndpoint(serv),
		DecodeGetAllRightsRequest,
		EncodeResponse,
		opts...,
	)

	getRightByNameHandler := kithttp.NewServer(
		MakeGetRightByNameEndpoint(serv),
		DecodeGetRightByNameRequest,
		EncodeResponse,
		opts...,
	)

	getRightByIdHandler := kithttp.NewServer(
		MakeGetRightByIdEndpoint(serv),
		DecodeGetRightByIdRequest,
		EncodeResponse,
		opts...,
	)

	deleteRightHandler := kithttp.NewServer(
		MakeDeleteRightEndpoint(serv),
		DecodeDeleteRightRequest,
		EncodeResponse,
		opts...,
	)
	return []*httphandlers.HTTPHandler{
		{
			Path:    "/api/v4/rights",
			Handler: addRightHandler,
			Methods: []string{"POST"},
		},
		{
			Path:    "/api/v4/rights/{rights_id}",
			Handler: editRightHandler,
			Methods: []string{"PUT"},
		},
		{
			Path:    "/api/v4/rights",
			Handler: getAllRightsHandler,
			Methods: []string{"GET"},
		},
		{
			Path:    "/api/v4/rights/{right_name}",
			Handler: getRightByNameHandler,
			Methods: []string{"GET"},
		},
		{
			Path:    "/api/v4/rights/id/{right_id}",
			Handler: getRightByIdHandler,
			Methods: []string{"GET"},
		},
		{
			Path:    "/api/v4/rights/{right_id}",
			Handler: deleteRightHandler,
			Methods: []string{"DELETE"},
		},
	}
}

type AddRightsRequest struct {
	Module string   `json:"module"`
	Action []string `json:"action"`
}

type AddRightsResponse struct {
	Err     error  `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

type EditRightRequest struct {
	ID     string   `json:"id"`
	Module string   `json:"module"`
	Action []string `json:"action"`
}

type EditRightResponse struct {
	Err     error  `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

type GetAllRightsRequest struct{}
type GetAllRightsResponse struct {
	Rights []domain.Right `json:"rights"`
	Err    error          `json:"error,omitempty"`
}

type GetRightByNameRequest struct {
	Module string `json:"module"`
}

type GetRightByNameResponse struct {
	Right *domain.Right `json:"right"`
	Err   error         `json:"error,omitempty"`
}

type GetRightByIdRequest struct {
	ID string `json:"id"`
}

type GetRightByIdResponse struct {
	Right *domain.Right `json:"right"`
	Err   error         `json:"error,omitempty"`
}

type DeleteRightRequest struct {
	ID string `json:"id"`
}

type DeleteRightResponse struct {
	Err     error  `json:"error,omitempty"`
	Message string `json:"message.omitempty"`
}

func MakeAddRightEndpoint(svc service.RightsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AddRightsRequest)
		err := svc.AddRights(ctx, req.Module, req.Action)
		return AddRightsResponse{Err: err, Message: "Права успешно созданы"}, err
	}
}

func MakeEditRightEndpoint(svc service.RightsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(EditRightRequest)
		err := svc.EditRight(ctx, req.ID, req.Module, req.Action)
		return EditRightResponse{Err: err, Message: "Права успешно обновлены"}, err
	}
}

func MakeGetAllRightsEndpoint(svc service.RightsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		rights, err := svc.GetAllRights(ctx)
		return GetAllRightsResponse{Rights: rights, Err: err}, err
	}
}

func MakeGetRightByNameEndpoint(svc service.RightsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetRightByNameRequest)
		right, err := svc.GetRightByName(ctx, req.Module)
		return GetRightByNameResponse{Right: right, Err: err}, err
	}
}

func MakeGetRightByIdEndpoint(svc service.RightsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetRightByIdRequest)
		right, err := svc.GetRightById(ctx, req.ID)
		return GetRightByIdResponse{Right: right, Err: err}, err
	}
}

func MakeDeleteRightEndpoint(svc service.RightsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteRightRequest)
		err := svc.DeleteRight(ctx, req.ID)
		return DeleteRightResponse{Err: err, Message: "Права успешно удалены"}, err
	}
}

func DecodeAddRightRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req AddRightsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func DecodeEditRightRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req EditRightRequest
	vars := mux.Vars(r)
	req.ID = vars["rights_id"]
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func DecodeGetAllRightsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return GetAllRightsRequest{}, nil
}

func DecodeGetRightByNameRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req GetRightByNameRequest
	vars := mux.Vars(r)
	req.Module = vars["right_name"]
	return req, nil
}

func DecodeGetRightByIdRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req GetRightByIdRequest
	vars := mux.Vars(r)
	req.ID = vars["right_id"]
	if req.ID == "" {
		return nil, errors.New("missing right_id")
	}
	return req, nil
}

func DecodeDeleteRightRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req DeleteRightRequest
	vars := mux.Vars(r)
	req.ID = vars["right_id"]
	return req, nil
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(response)
}

func NewAddRightsHandler(rightsService service.RightsService) http.Handler {
	return httptransport.NewServer(
		MakeAddRightEndpoint(rightsService),
		DecodeAddRightRequest,
		EncodeResponse,
	)
}
