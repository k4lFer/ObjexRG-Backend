package storage_models

import "github.com/google/uuid"

type ObjectShares struct {
	ID	uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ObjectId uuid.UUID `gorm:"type:uuid;not null;index"`
	OwnerId uuid.UUID `gorm:"type:uuid;not null;index"`
	TokenHash string `gorm:"type:text;not null;uniqueIndex"`
	ExpiresAt *int64 `gorm:"type:bigint;null"`
	Revoked bool `gorm:"not null;default:false"`
	DownloadCount int64 `gorm:"not null;default:0"`
	CreatedAt int64 `gorm:"type:bigint;not null"`

	Object Objects `gorm:"foreignKey:ObjectId;constraint:OnDelete:CASCADE;"`
}