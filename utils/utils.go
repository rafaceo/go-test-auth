package utils

import (
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	authServicePkg "github.com/rafaceo/go-test-auth/cmd/service"
	authHttp "github.com/rafaceo/go-test-auth/cmd/transport/https"
	userServiceFactory "github.com/rafaceo/go-test-auth/rights"
	userServicePkh "github.com/rafaceo/go-test-auth/rights/service"
	userHttp "github.com/rafaceo/go-test-auth/rights/transport/http"
)

func CreateHTTPRouting(authService authServicePkg.AuthService, userService userServicePkh.UserService, logger log.Logger, postgres *sqlx.DB) *mux.Router {
	userServiceFac := new(userServiceFactory.ServiceFactory).CreateUserService(logger, postgres)
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

	return r
}
