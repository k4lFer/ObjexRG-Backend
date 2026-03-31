package auth_models

import "github.com/google/uuid"

type UserCredentials struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserDataID uuid.UUID    `gorm:"type:uuid;not null"`
	Email      string    `gorm:"type:varchar(255);unique;not null"`
	Username   string    `gorm:"type:varchar(50);unique;not null"`
	Password   string    `gorm:"type:varchar(255);not null"`

	UserData UserData `gorm:"foreignKey:UserDataID;constraint:OnDelete:CASCADE"`
}