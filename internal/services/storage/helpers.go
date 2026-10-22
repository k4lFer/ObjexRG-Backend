package storage

import (
	"errors"
	"fmt"

	storage_repo "objexrg/internal/infrastructure/repositories/storage"
	storage_models "objexrg/internal/models/storage"

	"github.com/google/uuid"
)

func parseOwner(userID string) (uuid.UUID, error) {
	ownerUUID, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID: %w", err)
	}
	return ownerUUID, nil
}

// loadBucket resolves a bucket by name for the given owner, returning the
// owner UUID alongside it. It fails when the bucket is missing or not owned.
func loadBucket(repo *storage_repo.StorageRepository, userID, bucketName string) (*storage_models.Bucket, uuid.UUID, error) {
	ownerUUID, err := parseOwner(userID)
	if err != nil {
		return nil, uuid.Nil, err
	}

	bucket, err := repo.GetBucketByNameAndOwner(ownerUUID, bucketName)
	if err != nil {
		return nil, ownerUUID, fmt.Errorf("failed to load bucket: %w", err)
	}
	if bucket == nil {
		return nil, ownerUUID, errors.New("bucket not found or you don't have access")
	}
	return bucket, ownerUUID, nil
}
