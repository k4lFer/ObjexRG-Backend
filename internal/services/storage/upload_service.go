package storage

import (
	"fmt"
	"io"

	"objexrg-backend/internal/interfaces"
)

type UploadService struct {
	storage interfaces.StorageInterface
}

func NewUploadService(storage interfaces.StorageInterface) *UploadService {
	return &UploadService{storage: storage}
}

func (s *UploadService) Upload(bucket string, filename string, reader io.Reader) (string, error) {
	if bucket == "" {
		return "", fmt.Errorf("bucket es requerido")
	}
	if filename == "" {
		return "", fmt.Errorf("filename es requerido")
	}
	return s.storage.Save(bucket, reader, &filename)
}

func (s *UploadService) BucketExists(userID, bucketName string) bool {
	return s.storage.BucketExists(userID, bucketName)
}

func (s *UploadService) Download(bucket string, filename string) (io.ReadCloser, error) {
	return s.storage.Get(bucket, filename)
}
