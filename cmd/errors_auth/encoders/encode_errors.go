package encoders

import (
	"context"
	"encoding/json"
	"github.com/rafaceo/go-test-auth/cmd/errors_auth"
	"net/http"
)

func EncodeErrorAUTH(_ context.Context, err error, w http.ResponseWriter) {
	switch err {
	case errors_auth.Forbidden:
		w.WriteHeader(http.StatusForbidden)
	case errors_auth.BadRequestError:
		w.WriteHeader(http.StatusBadRequest)
	case errors_auth.UnauthorizedError:
		w.WriteHeader(http.StatusUnauthorized)
	case errors_auth.TooManyRequestError:
		w.WriteHeader(http.StatusTooManyRequests)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(err)
}
