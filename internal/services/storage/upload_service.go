package storage

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"time"

	storage_repo "objexrg-backend/internal/infrastructure/repositories/storage"
	"objexrg-backend/internal/interfaces"
	storage_models "objexrg-backend/internal/models/storage"

	"github.com/google/uuid"
)

// ErrObjectExists is returned when a file with the same name already exists
// in the target bucket.
var ErrObjectExists = errors.New("an object with that name already exists in this bucket")

type UploadService struct {
	storage interfaces.StorageInterface
	repo    *storage_repo.StorageRepository
}

func NewUploadService(storage interfaces.StorageInterface, repo *storage_repo.StorageRepository) *UploadService {
	return &UploadService{storage: storage, repo: repo}
}

func (s *UploadService) Upload(userID string, bucketName string, filename string, size int64, mimeType string, reader io.Reader) (string, error) {
	if bucketName == "" {
		return "", fmt.Errorf("bucket name is required")
	}
	if filename == "" {
		return "", fmt.Errorf("filename is required")
	}

	ownerUUID, err := uuid.Parse(userID)
	if err != nil {
		return "", fmt.Errorf("invalid user ID: %w", err)
	}

	bucket, err := s.repo.GetBucketByNameAndOwner(ownerUUID, bucketName)
	if err != nil || bucket == nil {
		return "", fmt.Errorf("bucket not found or error loading it")
	}

	existing, err := s.repo.GetObjectByFileName(bucket.ID, filename)
	if err != nil {
		return "", err
	}
	if existing != nil {
		return "", ErrObjectExists
	}

	bucketPath := userID + "/" + bucketName

	// Store the blob under a UUID-based key so the original name is only
	// metadata. This keeps the storage layer collision-free and opens the
	// door for content-addressed deduplication later.
	objectID := uuid.New()
	extension := filepath.Ext(filename)
	storageKey := objectID.String() + extension

	stored, err := s.storage.Save(bucketPath, reader, &storageKey)
	if err != nil {
		return "", err
	}

	obj := &storage_models.Objects{
		ID:         objectID,
		BucketID:   bucket.ID,
		OwnerID:    ownerUUID,
		Path:       bucketPath + "/" + storageKey,
		FileName:   filename,
		StorageKey: storageKey,
		SizeBytes:  size,
		MimeType:   mimeType,
		Extension:  extension,
		CreatedAt:  time.Now(),
	}

	if err := s.repo.CreateObject(obj); err != nil {
		// Roll back the blob so we don't leave orphaned files behind.
		_ = s.storage.DeleteObject(bucketPath, storageKey)
		return stored, fmt.Errorf("file saved but db record failed: %w", err)
	}

	return stored, nil
}

func (s *UploadService) BucketExists(userID, bucketName string) bool {
	return s.storage.BucketExists(userID, bucketName)
}

func (s *UploadService) Download(userID string, bucketName string, filename string) (io.ReadCloser, *storage_models.Objects, error) {
	ownerUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid user ID: %w", err)
	}

	bucket, err := s.repo.GetBucketByNameAndOwner(ownerUUID, bucketName)
	if err != nil || bucket == nil {
		return nil, nil, errors.New("bucket not found or you don't have access")
	}

	obj, err := s.repo.GetObjectByFileName(bucket.ID, filename)
	if err != nil {
		return nil, nil, err
	}
	if obj == nil {
		return nil, nil, errors.New("object not found")
	}

	bucketPath := userID + "/" + bucketName
	stream, err := s.storage.Get(bucketPath, obj.StorageName())
	if err != nil {
		return nil, nil, err
	}

	return stream, obj, nil
}
