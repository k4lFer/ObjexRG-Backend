package auth_input

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
