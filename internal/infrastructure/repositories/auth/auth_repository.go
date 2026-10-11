package authrepository

import (
	"errors"
	"time"

	"objexrg-backend/internal/models/auth"

	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) Create(user *auth_models.UserCredentials) error {
	return r.db.Create(user).Error
}

func (r *AuthRepository) GetByEmail(email string) (*auth_models.UserCredentials, error) {
	var user auth_models.UserCredentials
	err := r.db.Preload("UserData").Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) GetByUsername(username string) (*auth_models.UserCredentials, error) {
	var user auth_models.UserCredentials
	err := r.db.Preload("UserData").Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) CreateUserData(userData *auth_models.UserData) error {
	return r.db.Create(userData).Error
}

func (r *AuthRepository) RegisterUser(email, username, password, firstName, lastName string, birthDate time.Time) (string, error) {
	// Crear UserData primero
	userData := &auth_models.UserData{
		FirstName: firstName,
		LastName:  lastName,
		BirthDate: birthDate,
		CreatedAt: time.Now(),
	}
	if err := r.CreateUserData(userData); err != nil {
		return "", err
	}

	userCred := &auth_models.UserCredentials{
		UserDataID: userData.ID,
		Email:      email,
		Username:   username,
		Password:   password,
	}
	if err := r.Create(userCred); err != nil {
		return "", err
	}
	return userCred.ID.String(), nil
}

func (r *AuthRepository) GetByID(id string) (*auth_models.UserCredentials, error) {
	var user auth_models.UserCredentials
	err := r.db.Preload("UserData").Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) CreateRefreshToken(token *auth_models.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *AuthRepository) GetRefreshTokenByID(id string) (*auth_models.RefreshToken, error) {
	var token auth_models.RefreshToken
	err := r.db.Where("id = ?", id).First(&token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &token, nil
}

func (r *AuthRepository) RevokeRefreshToken(id string) error {
	return r.db.Model(&auth_models.RefreshToken{}).
		Where("id = ?", id).
		Update("revoked", true).Error
}
