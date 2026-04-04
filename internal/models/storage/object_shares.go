package storage_models

import (
	"time"

	"github.com/google/uuid"
)

type ObjectShares struct {
	ID	uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ObjectId uuid.UUID `gorm:"type:uuid;not null;index"`
	OwnerId uuid.UUID `gorm:"type:uuid;not null;index"`
	TokenHash string `gorm:"type:text;not null;uniqueIndex"`
	ExpiresAt *time.Time `gorm:"type:timestamp;null"`
	Revoked bool `gorm:"not null;default:false"`
	DownloadCount int64 `gorm:"not null;default:0"`
	CreatedAt time.Time `gorm:"type:timestamp;not null"`

	Object Objects `gorm:"foreignKey:ObjectId;constraint:OnDelete:CASCADE;"`
}