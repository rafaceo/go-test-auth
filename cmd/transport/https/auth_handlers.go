package https

import (
	kitlog "github.com/go-kit/kit/log"
	kittransport "github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"

	//"gitlab.fortebank.com/forte-market/apps/common-libs/encoders"
	//"gitlab.fortebank.com/forte-market/apps/common-libs/httphandlers"

	"github.com/rafaceo/go-test-auth/cmd/common-libs/encoders"
	"github.com/rafaceo/go-test-auth/cmd/common-libs/httphandlers"

	//"gitlab.fortebank.com/forte-market/apps/user-profile-api/src/authenticate/service"
	//"gitlab.fortebank.com/forte-market/apps/user-profile-api/src/authenticate/transport"

	"github.com/rafaceo/go-test-auth/cmd/service"
	"github.com/rafaceo/go-test-auth/cmd/transport"
)

func GetAuthHandlers(serv service.AuthService, logger kitlog.Logger) []*httphandlers.HTTPHandler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(kittransport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encoders.EncodeErrorJSON),
	}

	loginHandler := kithttp.NewServer(
		transport.MakeLoginEndpoint(serv),
		transport.DecodeLoginRequest,
		transport.EncodeLoginResponse,
		opts...,
	)

	loginHTTPHandler := &httphandlers.HTTPHandler{
		Path:    "/api/v4/auth/login",
		Handler: loginHandler,
		Methods: []string{"POST"},
	}

	return []*httphandlers.HTTPHandler{
		loginHTTPHandler,
	}
}
