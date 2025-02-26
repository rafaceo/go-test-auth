package domain

type Role struct {
	ID     int                 `json:"id,omitempty"`
	Name   string              `json:"role_name"`
	NameRu string              `json:"role_name_ru"`
	Notes  string              `json:"notes"`
	Rights map[string][]string `json:"rights"`
}
