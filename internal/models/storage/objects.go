package storage_models

import (
	"time"

	"github.com/google/uuid"
)

type Objects struct {
	ID	uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	BucketID uuid.UUID `gorm:"type:uuid;not null;index"`
	OwnerID uuid.UUID `gorm:"type:uuid;not null;index"`
	Path string `gorm:"type:text;not null"`
	FileName string `gorm:"type:text;not null"`
	SizeBytes int64 `gorm:"not null"`
	MimeType string `gorm:"type:text;not null"`
	Extension string `gorm:"type:text;not null"`
	StorageKey string `gorm:"type:text;null"`
	Hash string `gorm:"type:text;null"`
	IsPublic bool `gorm:"not null;default:false"`
	PublicKey string `gorm:"type:text;null"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null"`
	DeletedAt *time.Time `gorm:"type:timestamptz;null"`

	Bucket Bucket `gorm:"foreignKey:BucketID;constraint:OnDelete:CASCADE;"`

}

// StorageName returns the on-disk/S3 key of the blob. New objects are stored
// under a UUID-based key; legacy rows fall back to the original file name.
func (o *Objects) StorageName() string {
	if o.StorageKey != "" {
		return o.StorageKey
	}
	return o.FileName
}
