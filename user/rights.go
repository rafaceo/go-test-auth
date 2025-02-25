package user

import (
	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	"github.com/rafaceo/go-test-auth/common-libs/instrumenting"
	"github.com/rafaceo/go-test-auth/user/middleware"
	"github.com/rafaceo/go-test-auth/user/repository/postgres"
	"github.com/rafaceo/go-test-auth/user/service"
)

type ServiceFactory struct{}

func (sf *ServiceFactory) CreateUserService(logger log.Logger, postgresClient *sqlx.DB) service.UserService {
	userRepo := postgres.NewUserRepository(postgresClient)

	userServ := service.NewUserService(userRepo)
	userServ = middleware.NewLoggingMiddleware(log.With(logger, "component", "users"), userServ)

	counter, duration, counterError := instrumenting.GetMetricsBySubsystem("user_service")
	userServ = middleware.NewInstrumentingMiddleware(counter, duration, counterError, userServ)

	return userServ
}
