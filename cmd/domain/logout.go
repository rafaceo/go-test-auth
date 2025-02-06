package domain

type LogoutRequest struct {
	RefreshToken string `json:"refreshToken"`
}
