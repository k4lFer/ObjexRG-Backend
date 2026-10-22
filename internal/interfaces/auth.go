package interfaces

import (
	"objexrg/internal/models/auth"
	"time"
)

type AuthInterface interface {
	Create(user *auth_models.UserCredentials) error
	GetByEmail(email string) (*auth_models.UserCredentials, error)
	GetByUsername(username string) (*auth_models.UserCredentials, error)
	RegisterUser(email, username, password, firstName, lastName string, birthDate time.Time) (string, error)
	GetByID(id string) (*auth_models.UserCredentials, error)

	CreateRefreshToken(token *auth_models.RefreshToken) error
	GetRefreshTokenByID(id string) (*auth_models.RefreshToken, error)
	RevokeRefreshToken(id string) error
}
