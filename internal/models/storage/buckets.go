package storage_models

import (
	"time"
	
	"github.com/google/uuid"
)

type Bucket struct {
    ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    OwnerID uuid.UUID `gorm:"type:uuid;not null;index:idx_owner_name,unique"`
    Name string `gorm:"type:varchar(255);not null;index:idx_owner_name,unique"`
    Status string `gorm:"type:varchar(50);not null;default:'active'"`
    CreatedAt time.Time `gorm:"type:timestamptz;not null"`
}