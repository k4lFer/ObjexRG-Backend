package auth_response

type LoginResponse struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}