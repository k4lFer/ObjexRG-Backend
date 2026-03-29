package user_response

import "time"

type MeResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	BirthDate time.Time    `json:"birth_date"`
	CreatedAt time.Time `json:"created_at"`
}