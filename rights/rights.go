package rights

import (
	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	"github.com/rafaceo/go-test-auth/common-libs/instrumenting"
	"github.com/rafaceo/go-test-auth/rights/middleware"
	"github.com/rafaceo/go-test-auth/rights/repository/postgres"
	"github.com/rafaceo/go-test-auth/rights/service"
)

// NewRightsRouter создает роутер с хендлерами и middleware
func CreateRightsService(logger log.Logger, postgresClient *sqlx.DB) service.RightsService {
	rightsRepo := postgres.NewPostgresRightsRepository(postgresClient)

	rightsServ := service.NewRightsService(rightsRepo)
	rightsServ = middleware.NewLoggingMiddleware(log.With(logger, "component", "rights"), rightsServ)
	counter, duration, counterError := instrumenting.GetMetricsBySubsystem("user_profile_service")
	rightsServ = middleware.NewInstrumentingMiddleware(counter, duration, counterError, rightsServ)

	return rightsServ
}
