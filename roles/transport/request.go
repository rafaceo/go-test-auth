package transport

type AddRoleRequest struct {
	RoleName   string              `json:"role_name"`
	RoleNameRu string              `json:"role_name_ru"`
	Notes      string              `json:"notes"`
	Rights     map[string][]string `json:"rights"`
}

type EditRoleRequest struct {
	RoleID     int                 `json:"-"`
	RoleName   string              `json:"role_name"`
	RoleNameRu string              `json:"role_name_ru"`
	Notes      string              `json:"notes"`
	Rights     map[string][]string `json:"rights"`
}

type GetRoleRequest struct {
}

type GetRoleRightsRequest struct {
	RoleID int `json:"-"`
}

type DeleteRoleRequest struct {
	RoleID int `json:"role_id"`
}

type AssignRoleToUserRequest struct {
	RoleID int  `json:"role_id"`
	UserID int  `json:"user_id"`
	Merge  bool `json:"merge"`
}
