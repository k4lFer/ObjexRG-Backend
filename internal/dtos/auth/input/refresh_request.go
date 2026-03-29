package auth_input

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}