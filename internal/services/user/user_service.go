package user_service

import (
	"objexrg-backend/internal/interfaces"
	models "objexrg-backend/internal/models/auth"

	"gorm.io/gorm"
)

type UserService struct {
	user interfaces.AuthInterface
	db   *gorm.DB
}

func NewUserService(user interfaces.AuthInterface, db *gorm.DB) *UserService {
	return &UserService{user: user, db: db}
}

func (s *UserService) GetById(id string) (*models.UserCredentials, error) {
	return s.user.GetByID(id)
}