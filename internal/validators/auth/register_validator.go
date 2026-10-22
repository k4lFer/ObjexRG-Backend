package authvalidators

import (
	"errors"
	"net/http"
	dtos "objexrg/internal/dtos/auth/input"
	"strings"
	"time"

	"objexrg/internal/interfaces"
	"objexrg/internal/validators"

	"github.com/go-playground/validator/v10"
)

type RegisterValidator struct {
	authRepo interfaces.AuthInterface
	validate *validator.Validate
}

func NewRegisterValidator(authRepo interfaces.AuthInterface) *RegisterValidator {
	v := validator.New()
	v.RegisterValidation("birthdate", func(fl validator.FieldLevel) bool {
		dateStr := fl.Field().String()
		_, err := time.Parse("2006-01-02", dateStr)
		return err == nil
	})
	return &RegisterValidator{authRepo: authRepo, validate: v}
}

func (v *RegisterValidator) Validate(input interface{}) validators.ValidationResult {
	req, ok := input.(*dtos.RegisterRequest)
	if !ok {
		return validators.ValidationResult{
			IsValid:    false,
			StatusCode: http.StatusInternalServerError,
			Message:    "Invalid input type",
		}
	}

	// Validar formato con validator
	if err := v.validate.Struct(req); err != nil {
		return validators.ValidationResult{
			IsValid:    false,
			StatusCode: http.StatusUnprocessableEntity,
			Message:    err.Error(),
		}
	}

	// Validar existencia (duplicados)
	if _, err := v.authRepo.GetByEmail(req.Email); err == nil {
		return validators.ValidationResult{
			IsValid:    false,
			StatusCode: http.StatusConflict,
			Message:    "Email already exists",
		}
	}
	if _, err := v.authRepo.GetByUsername(req.Username); !errors.Is(err, errors.New("user not found")) && !strings.Contains(err.Error(), "user not found") {
		// Si no es "user not found", significa que existe
		if err == nil {
			return validators.ValidationResult{
				IsValid:    false,
				StatusCode: http.StatusConflict,
				Message:    "Username already exists",
			}
		}
	}

	return validators.ValidationResult{
		IsValid:    true,
		StatusCode: http.StatusOK,
		Message:    "",
	}
}
