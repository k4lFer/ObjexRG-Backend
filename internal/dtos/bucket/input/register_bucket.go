package bucket_input

// RegisterRequest represents the request to register a new bucket
type RegisterRequest struct {
	Name	string	`json:"name" validate:"required,min=3,max=50"`
}
