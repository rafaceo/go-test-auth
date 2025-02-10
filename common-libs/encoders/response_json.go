package encoders

import (
	"context"
	"encoding/json"
	"github.com/rafaceo/go-test-auth/common-libs/interfaces"
	"net/http"
)

func EncodeResponseJSON(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(interfaces.Errorer); ok && e.Error() != nil {
		EncodeErrorJSON(ctx, e.Error(), w)

		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	return json.NewEncoder(w).Encode(response)
}
