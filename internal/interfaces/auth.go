package interfaces

import (
	models "objexrg-backend/internal/models/auth"
	"time"
)

type Auth interface {
	Create(user *models.UserCredentials) error
	GetByEmail(email string) (*models.UserCredentials, error)
	GetByUsername(username string) (*models.UserCredentials, error)
	RegisterUser(email, username, password, firstName, lastName string, birthDate time.Time) (string, error)
	GetByID(id string) (*models.UserCredentials, error)
}
