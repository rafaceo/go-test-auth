package domain

// SecretKey для подписи JWT (должен быть в .env)
var SecretKey = []byte("")

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type RefreshResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
