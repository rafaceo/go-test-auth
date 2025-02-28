package roles

import (
	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	"gitlab.fortebank.com/forte-market/apps/common-libs/instrumenting"
	"gitlab.fortebank.com/forte-market/apps/user-profile-api/src/roles/middleware"
	"gitlab.fortebank.com/forte-market/apps/user-profile-api/src/roles/repository/postgres"
	"gitlab.fortebank.com/forte-market/apps/user-profile-api/src/roles/service"
)

// NewRolesRouter создает роутер с хендлерами и middleware
func CreateRolesService(logger log.Logger, postgresClient *sqlx.DB) service.RoleService {
	rolesRepo := postgres.NewRoleRepository(postgresClient)
	roleService := service.NewRoleService(rolesRepo)
	roleService = middleware.NewLoggingMiddleware(log.With(logger, "component", "roles"), roleService)
	counter, duration, counterError := instrumenting.GetMetricsBySubsystem("user_profile_service")
	roleService = middleware.NewInstrumentingMiddleware(counter, duration, counterError, roleService)

	return roleService
}
