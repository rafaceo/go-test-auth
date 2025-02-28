package user_contexts

import (
	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	"github.com/rafaceo/go-test-auth/common-libs/instrumenting"
	"github.com/rafaceo/go-test-auth/user_contexts/middleware"
	"github.com/rafaceo/go-test-auth/user_contexts/repository/postgres"
	"github.com/rafaceo/go-test-auth/user_contexts/service"
)

func CreateUserContextRouter(logger log.Logger, postgresClient *sqlx.DB) service.UserContextService {
	userCtxRepo := postgres.NewUserContextRepository(postgresClient)
	userCtxService := service.NewUserContextService(userCtxRepo)
	userCtxService = middleware.NewLoggingMiddleware(log.With(logger, "component", "user_contexts"), userCtxService)
	counter, duration, counterError := instrumenting.GetMetricsBySubsystem("user_context_service")
	userCtxService = middleware.NewInstrumentingMiddleware(counter, duration, counterError, userCtxService)

	return userCtxService
}
