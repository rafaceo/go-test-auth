package encoders

import (
	"context"
	"encoding/json"
	errors2 "github.com/rafaceo/go-test-auth/cmd/common-libs/errors"
	"net/http"
)

func EncodeErrorJSON(_ context.Context, err error, w http.ResponseWriter) {
	switch err {
	case errors2.OK:
		w.WriteHeader(http.StatusOK)
	case errors2.Conflict, errors2.CsvError, errors2.CassandraReadError:
		w.WriteHeader(http.StatusConflict)
	case errors2.NotFound:
		w.WriteHeader(http.StatusNotFound)
	case errors2.AccessDenied:
		w.WriteHeader(http.StatusForbidden)
	case errors2.ElasticConnectError, errors2.S3ConnectError, errors2.CassandraConnectError, errors2.RabbitMQConnectError:
		w.WriteHeader(http.StatusServiceUnavailable)
	case errors2.ContentNotFound:
		w.WriteHeader(http.StatusNoContent)
	case errors2.DeserializeBug:
		w.WriteHeader(http.StatusUnsupportedMediaType)
	case errors2.InvalidCharacter:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(err)
}
