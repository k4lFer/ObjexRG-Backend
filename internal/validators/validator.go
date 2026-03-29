package validators

type ValidationResult struct {
	IsValid    bool
	StatusCode int
	Message    string
}

type Validator interface {
	Validate(input interface{}) ValidationResult
}
