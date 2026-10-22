package blob

import (
	"context"
	"fmt"

	"objexrg/internal/config"
	"objexrg/internal/interfaces"
)

const (
	DriverLocal = "local"
	DriverS3    = "s3"
)

// NewBlobStore builds the BlobStore selected by STORAGE_DRIVER. For the S3
// driver the backing bucket is created on startup when missing.
func NewBlobStore(ctx context.Context, cfg *config.Config) (interfaces.BlobStore, error) {
	switch cfg.StorageDriver {
	case DriverS3:
		store, err := NewS3BlobStore(S3Options{
			Endpoint:       cfg.S3Endpoint,
			PublicEndpoint: cfg.S3PublicEndpoint,
			Region:         cfg.S3Region,
			AccessKey:      cfg.S3AccessKey,
			SecretKey:      cfg.S3SecretKey,
			Bucket:         cfg.S3Bucket,
			UseSSL:         cfg.S3UseSSL,
			ForcePathStyle: cfg.S3ForcePathStyle,
		})
		if err != nil {
			return nil, err
		}
		if err := store.EnsureBucket(ctx); err != nil {
			return nil, err
		}
		return store, nil

	case DriverLocal, "":
		return NewLocalBlobStore(cfg.BasePath), nil

	default:
		return nil, fmt.Errorf("unknown STORAGE_DRIVER %q (expected %q or %q)", cfg.StorageDriver, DriverLocal, DriverS3)
	}
}
