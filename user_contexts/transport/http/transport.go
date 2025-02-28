package http

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	kittransport "github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/rafaceo/go-test-auth/cmd/errors_auth/encoders"
	"github.com/rafaceo/go-test-auth/common-libs/httphandlers"
	"github.com/rafaceo/go-test-auth/user_contexts/service"
	"github.com/rafaceo/go-test-auth/user_contexts/transport"
)

func GetUserContextHandlers(serv service.UserContextService, logger kitlog.Logger) []*httphandlers.HTTPHandler {
	opts := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(kittransport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encoders.EncodeErrorAUTH),
	}

	addUserContextHandler := httptransport.NewServer(
		MakeAddUserContextEndpoint(serv),
		DecodeAddUserContextRequest,
		EncodeResponse,
		opts...,
	)

	editUserContextHandler := httptransport.NewServer(
		MakeEditUserContextEndpoint(serv),
		DecodeEditUserContextRequest,
		EncodeResponse,
		opts...,
	)

	getUserContextsHandler := httptransport.NewServer(
		MakeGetUserContextsEndpoint(serv),
		DecodeGetUserContextsRequest,
		EncodeResponse,
		opts...,
	)

	deleteUserContextHandler := httptransport.NewServer(
		MakeDeleteUserContextEndpoint(serv),
		DecodeDeleteUserContextRequest,
		EncodeResponse,
		opts...,
	)

	deleteAllUserContextsHandler := httptransport.NewServer(
		MakeDeleteAllUserContextsEndpoint(serv),
		DecodeDeleteAllUserContextsRequest,
		EncodeResponse,
		opts...,
	)

	return []*httphandlers.HTTPHandler{
		{
			Path:    "/api/v4/users/{user_id}/contexts",
			Handler: addUserContextHandler,
			Methods: []string{"POST"},
		},
		{
			Path:    "/api/v4/users/{user_id}/contexts",
			Handler: editUserContextHandler,
			Methods: []string{"PUT"},
		},
		{
			Path:    "/api/v4/users/{user_id}/contexts",
			Handler: getUserContextsHandler,
			Methods: []string{"GET"},
		},
		{
			Path:    "/api/v4/users/{user_id}/contexts/{merchant_id}",
			Handler: deleteUserContextHandler,
			Methods: []string{"DELETE"},
		},
		{
			Path:    "/api/v4/users/{user_id}/contexts",
			Handler: deleteAllUserContextsHandler,
			Methods: []string{"DELETE"},
		},
	}
}

// Эндпоинт для добавления контекста
func MakeAddUserContextEndpoint(svc service.UserContextService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(transport.AddUserContextRequest)
		err := svc.AddUserContext(ctx, req.UserID, req.MerchantID, req.Global)
		if err != nil {
			return transport.AddUserContextResponse{Error: err.Error()}, nil
		}
		return transport.AddUserContextResponse{Message: "Context added"}, nil
	}
}

// Эндпоинт для редактирования контекста
func MakeEditUserContextEndpoint(svc service.UserContextService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(transport.EditUserContextRequest)
		err := svc.EditUserContext(ctx, req.UserID, req.Global)
		if err != nil {
			return transport.EditUserContextResponse{Error: err.Error()}, nil
		}
		return transport.EditUserContextResponse{}, nil
	}
}

func MakeGetUserContextsEndpoint(svc service.UserContextService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(transport.GetUserContextRequest)
		if !ok {
			return transport.GetUserContextResponse{Error: "invalid request format"}, nil
		}

		contexts, err := svc.GetUserContexts(ctx, req.UserID)
		if err != nil {
			return transport.GetUserContextResponse{Error: err.Error()}, nil
		}

		return transport.GetUserContextResponse{UserContext: contexts}, nil
	}
}

// Эндпоинт для удаления одного контекста пользователя
func MakeDeleteUserContextEndpoint(svc service.UserContextService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(transport.DeleteUserContextRequest)
		err := svc.DeleteUserContext(ctx, req.UserID, req.MerchantID)
		if err != nil {
			return transport.DeleteUserContextResponse{Error: err.Error()}, nil
		}
		return transport.DeleteUserContextResponse{}, nil
	}
}

// Эндпоинт для удаления всех контекстов пользователя
func MakeDeleteAllUserContextsEndpoint(svc service.UserContextService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(transport.DeleteAllUserContextsRequest)
		err := svc.DeleteAllUserContexts(ctx, req.UserID)
		if err != nil {
			return transport.DeleteAllUserContextsResponse{Error: err.Error()}, nil
		}
		return transport.DeleteAllUserContextsResponse{}, nil
	}
}

// Декодирование запроса на добавление контекста
func DecodeAddUserContextRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req transport.AddUserContextRequest

	// Декодируем JSON в req
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	defer r.Body.Close()

	// Получаем user_id из URL
	vars := mux.Vars(r)
	userIDStr, ok := vars["user_id"]
	if !ok {
		return nil, errors.New("missing user_id in URL")
	}

	// Конвертируем строку в uuid.UUID
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user_id format")
	}

	// Заполняем user_id в req
	req.UserID = userID

	// Проверяем, что merchant_id не пустой
	if req.MerchantID == "" {
		return nil, errors.New("merchant_id cannot be empty")
	}

	return req, nil
}

// Декодирование запроса на редактирование контекста
func DecodeEditUserContextRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req transport.EditUserContextRequest

	// Декодируем JSON в req
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	defer r.Body.Close()

	// Получаем user_id из URL
	vars := mux.Vars(r)
	userIDStr, ok := vars["user_id"]
	if !ok {
		return nil, errors.New("missing user_id in URL")
	}

	// Конвертируем строку в uuid.UUID
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user_id format")
	}

	// Заполняем user_id в req
	req.UserID = userID

	return req, nil
}

// Декодирование запроса на получение контекстов пользователя
func DecodeGetUserContextsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	userIDStr, ok := vars["user_id"]
	if !ok {
		return nil, errors.New("missing user_id in URL")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user_id format")
	}

	return transport.GetUserContextRequest{UserID: userID}, nil
}

// Декодирование запроса на удаление одного контекста пользователя
func DecodeDeleteUserContextRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	userIDStr, ok := vars["user_id"]
	if !ok {
		return nil, errors.New("missing user_id in URL")
	}

	merchantID, ok := vars["merchant_id"]
	if !ok {
		return nil, errors.New("missing merchant_id in URL")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user_id format")
	}

	return transport.DeleteUserContextRequest{
		UserID:     userID,
		MerchantID: merchantID,
	}, nil
}

// Декодирование запроса на удаление всех контекстов пользователя
func DecodeDeleteAllUserContextsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	userIDStr, ok := vars["user_id"]
	if !ok {
		return nil, errors.New("missing user_id in URL")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user_id format")
	}

	return transport.DeleteAllUserContextsRequest{UserID: userID}, nil
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(response)
}
