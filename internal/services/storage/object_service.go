package storage

import (
	"errors"
	"fmt"
	"io"

	storage_repo "objexrg-backend/internal/infrastructure/repositories/storage"
	"objexrg-backend/internal/interfaces"
	storage_models "objexrg-backend/internal/models/storage"

	"github.com/google/uuid"
)

type ObjectService struct {
	storage interfaces.StorageInterface
	repo    *storage_repo.StorageRepository
}

func NewObjectService(storage interfaces.StorageInterface, repo *storage_repo.StorageRepository) *ObjectService {
	return &ObjectService{storage: storage, repo: repo}
}

// SetVisibility toggles whether an object can be read through the public route.
func (s *ObjectService) SetVisibility(userID, bucketName, filename string, isPublic bool) (*storage_models.Objects, error) {
	ownerUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	bucket, err := s.repo.GetBucketByNameAndOwner(ownerUUID, bucketName)
	if err != nil || bucket == nil {
		return nil, fmt.Errorf("bucket not found or you don't have access")
	}

	if err := s.repo.UpdateObjectVisibility(ownerUUID, bucket.ID, filename, isPublic); err != nil {
		return nil, err
	}

	obj, err := s.repo.GetObjectByFileName(bucket.ID, filename)
	if err != nil || obj == nil {
		return nil, fmt.Errorf("object not found")
	}

	return obj, nil
}

// DeleteObject removes the file from disk and its database record.
func (s *ObjectService) DeleteObject(userID, bucketName, filename string) error {
	ownerUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	bucket, err := s.repo.GetBucketByNameAndOwner(ownerUUID, bucketName)
	if err != nil || bucket == nil {
		return fmt.Errorf("bucket not found or you don't have access")
	}

	obj, err := s.repo.GetObjectByFileName(bucket.ID, filename)
	if err != nil {
		return err
	}
	if obj == nil {
		return fmt.Errorf("object not found")
	}

	bucketPath := userID + "/" + bucketName
	if err := s.storage.DeleteObject(bucketPath, filename); err != nil {
		return fmt.Errorf("failed to remove file from storage: %w", err)
	}

	return s.repo.DeleteObject(ownerUUID, bucket.ID, filename)
}

// GetPublic streams an object that has been explicitly marked as public.
func (s *ObjectService) GetPublic(objectID string) (io.ReadCloser, *storage_models.Objects, error) {
	id, err := uuid.Parse(objectID)
	if err != nil {
		return nil, nil, errors.New("invalid object id")
	}

	obj, err := s.repo.GetObjectByID(id)
	if err != nil {
		return nil, nil, err
	}
	if obj == nil || !obj.IsPublic {
		return nil, nil, errors.New("object not found")
	}

	bucketPath := obj.OwnerID.String() + "/" + obj.Bucket.Name
	stream, err := s.storage.Get(bucketPath, obj.FileName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to retrieve file from storage: %w", err)
	}

	return stream, obj, nil
}
