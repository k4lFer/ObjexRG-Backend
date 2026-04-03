package storage

import (
	"fmt"
	"io"
	"path/filepath"
	"time"

	storage_repo "objexrg-backend/internal/infrastructure/repositories/storage"
	"objexrg-backend/internal/interfaces"
	storage_models "objexrg-backend/internal/models/storage"

	"github.com/google/uuid"
)

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

	bucketPath := userID + "/" + bucketName
	stored, err := s.storage.Save(bucketPath, reader, &filename)
	if err != nil {
		return "", err
	}

	extension := filepath.Ext(filename)

	obj := &storage_models.Objects{
		BucketID:  bucket.ID,
		OwnerID:   ownerUUID,
		Path:      bucketPath + "/" + filename, // internal storage path
		FileName:  filename,
		SizeBytes: size,
		MimeType:  mimeType,
		Extension: extension,
		CreatedAt: time.Now(),
	}

	if err := s.repo.CreateObject(obj); err != nil {
		return stored, fmt.Errorf("file saved but db record failed: %w", err)
	}

	return stored, nil
}

func (s *UploadService) BucketExists(userID, bucketName string) bool {
	return s.storage.BucketExists(userID, bucketName)
}

func (s *UploadService) Download(bucket string, filename string) (io.ReadCloser, error) {
	return s.storage.Get(bucket, filename)
}

