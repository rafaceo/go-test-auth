package domain

var SecretKey = []byte("sjCN/AMyXYwufn/hF3FcikY7SORti8wYX+lp8dh+V+I=")

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type RefreshResponse struct {
	AccessToken string `json:"accessToken"`
}
