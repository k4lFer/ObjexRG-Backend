package authvalidators

import (
	"net/http"
	dtos "objexrg/internal/dtos/auth/input"

	"objexrg/internal/interfaces"
	"objexrg/internal/validators"

	"github.com/go-playground/validator/v10"
)

type LoginValidator struct {
	authRepo interfaces.AuthInterface
	validate *validator.Validate
}

func NewLoginValidator(authRepo interfaces.AuthInterface) *LoginValidator {
	v := validator.New()
	return &LoginValidator{authRepo: authRepo, validate: v}
}

func (v *LoginValidator) Validate(input interface{}) validators.ValidationResult {
	req, ok := input.(*dtos.LoginRequest)
	if !ok {
		return validators.ValidationResult{
			IsValid:    false,
			StatusCode: http.StatusInternalServerError,
			Message:    "Invalid input type",
		}
	}

	// Validar formato
	if err := v.validate.Struct(req); err != nil {
		return validators.ValidationResult{
			IsValid:    false,
			StatusCode: http.StatusUnprocessableEntity,
			Message:    err.Error(),
		}
	}

	// Para login, no validar existencia aquí, ya que eso se hace en servicio para credenciales inválidas
	return validators.ValidationResult{
		IsValid:    true,
		StatusCode: http.StatusOK,
		Message:    "",
	}
}
