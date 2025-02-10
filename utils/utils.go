package utils

import (
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	authServicePkg "github.com/rafaceo/go-test-auth/cmd/service"
	authHttp "github.com/rafaceo/go-test-auth/cmd/transport/https"
)

func CreateHTTPRouting(authService authServicePkg.AuthService, logger log.Logger, postgres *sqlx.DB) *mux.Router {

	r := mux.NewRouter()

	authHTTPHandlers := authHttp.GetAuthHandlers(authService, logger)
	if len(authHTTPHandlers) > 0 {
		for _, authHTTPHandler := range authHTTPHandlers {
			r.Handle(authHTTPHandler.Path, authHTTPHandler.Handler).Methods(authHTTPHandler.Methods...)
		}
	}

	return r
}
