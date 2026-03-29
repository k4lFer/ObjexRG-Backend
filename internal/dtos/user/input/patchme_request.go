package user_input

type PatchMeRequest struct {
	Username  *string `json:"username,omitempty"`
	Email     *string `json:"email,omitempty"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	BirthDate *string `json:"birth_date,omitempty"` // Formato "YYYY-MM-DD"
}