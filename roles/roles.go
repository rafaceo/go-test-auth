package roles

import (
	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	"github.com/rafaceo/go-test-auth/common-libs/instrumenting"
	"github.com/rafaceo/go-test-auth/roles/middleware"
	"github.com/rafaceo/go-test-auth/roles/repository/postgres"
	"github.com/rafaceo/go-test-auth/roles/service"
)

type ServiceFactory struct{}

func (sf *ServiceFactory) CreateRolesService(logger log.Logger, postgresClient *sqlx.DB) service.RoleService {
	rolesRepo := postgres.NewRoleRepository(postgresClient)
	roleService := service.NewRoleService(rolesRepo)
	roleService = middleware.NewLoggingMiddleware(log.With(logger, "component", "roles"), roleService)
	counter, duration, counterError := instrumenting.GetMetricsBySubsystem("user_profile_service")
	roleService = middleware.NewInstrumentingMiddleware(counter, duration, counterError, roleService)

	return roleService
}
