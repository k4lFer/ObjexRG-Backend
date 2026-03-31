package auth_models

import (
	"time"

	"github.com/google/uuid"
)

type UserData struct {
	ID uuid.UUID		`gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	FirstName string     `gorm:"type:varchar(255);not null"`
	LastName  string     `gorm:"type:varchar(255);not null"`
	BirthDate time.Time  `gorm:"type:date;not null"`
	CreatedAt time.Time  `gorm:"type:timestamptz;not null"`
	UpdatedAt *time.Time `gorm:"type:timestamptz"`
}