package auth_models

import (
	"time"

	"github.com/google/uuid"
)

// RefreshToken stores a single issued refresh token (identified by its JTI)
// so it can be rotated and revoked on logout.
type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	ExpiresAt time.Time `gorm:"type:timestamptz;not null"`
	Revoked   bool      `gorm:"not null;default:false"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null"`
}
