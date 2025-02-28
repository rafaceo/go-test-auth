package utils

import (
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	authServicePkg "github.com/rafaceo/go-test-auth/cmd/service"
	authHttp "github.com/rafaceo/go-test-auth/cmd/transport/https"
	rightsServicePkg "github.com/rafaceo/go-test-auth/rights/service"
	rightsHttp "github.com/rafaceo/go-test-auth/rights/transport/http"
	rolesServiceFactory "github.com/rafaceo/go-test-auth/roles"
	rolesHttp "github.com/rafaceo/go-test-auth/roles/transport/http"
	userServiceFactory "github.com/rafaceo/go-test-auth/user"
	userHttp "github.com/rafaceo/go-test-auth/user/transport/http"
)

func CreateHTTPRouting(authService authServicePkg.AuthService, rightsService rightsServicePkg.RightsService, logger log.Logger, postgres *sqlx.DB) *mux.Router {
	userServiceFac := new(userServiceFactory.ServiceFactory).CreateUserService(logger, postgres)
	rolesServiceFac := new(rolesServiceFactory.ServiceFactory).CreateRolesService(logger, postgres)
	r := mux.NewRouter()
	userHTTPHandlers := userHttp.GetUserHandler(userServiceFac, logger)
	if len(userHTTPHandlers) > 0 {
		for _, userHTTPHandler := range userHTTPHandlers {
			r.Handle(userHTTPHandler.Path, userHTTPHandler.Handler).Methods(userHTTPHandler.Methods...)
		}
	}

	authHTTPHandlers := authHttp.GetAuthHandlers(authService, logger)
	if len(authHTTPHandlers) > 0 {
		for _, authHTTPHandler := range authHTTPHandlers {
			r.Handle(authHTTPHandler.Path, authHTTPHandler.Handler).Methods(authHTTPHandler.Methods...)
		}
	}

	rolesHTTPHandlers := rolesHttp.GetRoleHandlers(rolesServiceFac, logger)
	if len(rolesHTTPHandlers) > 0 {
		for _, rolesHTTPHandler := range rolesHTTPHandlers {
			r.Handle(rolesHTTPHandler.Path, rolesHTTPHandler.Handler).Methods(rolesHTTPHandler.Methods...)
		}
	}

	rightsHTTPHandlers := rightsHttp.GetRightHandlers(rightsService, logger)
	if len(rightsHTTPHandlers) > 0 {
		for _, rightsHTTPHandler := range rightsHTTPHandlers {
			r.Handle(rightsHTTPHandler.Path, rightsHTTPHandler.Handler).Methods(rightsHTTPHandler.Methods...)
		}
	}

	return r
}
