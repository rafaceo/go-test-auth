package domain

import "encoding/json"

type User struct {
	ID      uint            `json:"id"`
	Name    string          `json:"name"`
	Context string          `json:"context"`
	Rights  json.RawMessage `json:"rights"`
}
