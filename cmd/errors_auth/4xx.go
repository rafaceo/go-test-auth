package errors_auth

import "fmt"

// 4XX
var (
	Forbidden           = &ArgError{ArgErrorSystemMarket, 403, "forbidden", "forbidden"}
	BadRequestError     = &ArgError{ArgErrorSystemMarket, 400, "Missing required fields: phone, password", "Required fields missing: phone, password"}
	UnauthorizedError   = &ArgError{ArgErrorSystemMarket, 401, "Incorrect login or password", "Unauthorized: Incorrect login or password"}
	TooManyRequestError = &ArgError{ArgErrorSystemMarket, 429, "too many requests", "Too many requests: 1 min"}
)

const (
	ArgErrorSystemMarket = "MARKET"
)

// ///////////////////////////////
// Error structure & methods
// ///////////////////////////////
type ArgError struct {
	System           string `json:"system"`
	Status           int    `json:"status"`
	Message          string `json:"message"`
	DeveloperMessage string `json:"developerMessage"`
}

func (e *ArgError) Error() string {
	return fmt.Sprintf("%d %s", e.Status, e.DeveloperMessage)
}

func (e *ArgError) SetDevMessage(developMessage string) *ArgError {
	e.DeveloperMessage = developMessage

	return e
}
