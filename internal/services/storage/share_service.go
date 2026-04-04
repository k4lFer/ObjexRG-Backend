package storage

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"time"

	storage_repo "objexrg-backend/internal/infrastructure/repositories/storage"
	"objexrg-backend/internal/interfaces"
	storage_models "objexrg-backend/internal/models/storage"

	"github.com/google/uuid"
)

type ShareService struct {
	storage interfaces.StorageInterface
	repo    *storage_repo.StorageRepository
}

func NewShareService(storage interfaces.StorageInterface, repo *storage_repo.StorageRepository) *ShareService {
	return &ShareService{storage: storage, repo: repo}
}

func generateToken(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b), nil
}

func (s *ShareService) GenerateShareLink(userID string, bucketName string, filename string, expiresAt *time.Time) (string, error) {
	if bucketName == "" || filename == "" {
		return "", fmt.Errorf("bucket name and filename are required")
	}

	ownerUUID, err := uuid.Parse(userID)
	if err != nil {
		return "", fmt.Errorf("invalid user ID: %w", err)
	}

	bucket, err := s.repo.GetBucketByNameAndOwner(ownerUUID, bucketName)
	if err != nil || bucket == nil {
		return "", fmt.Errorf("bucket not found or you don't have access")
	}

	obj, err := s.repo.GetObjectByFileName(bucket.ID, filename)
	if err != nil || obj == nil {
		return "", fmt.Errorf("object not found in this bucket")
	}

	token, err := generateToken(24) // 24 bytes = ~32 base64 chars
	if err != nil {
		return "", fmt.Errorf("failed to generate secure token: %w", err)
	}

	share := &storage_models.ObjectShares{
		ObjectId:  obj.ID,
		OwnerId:   ownerUUID,
		TokenHash: token,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}

	if err := s.repo.CreateObjectShare(share); err != nil {
		return "", fmt.Errorf("failed to create db share record: %w", err)
	}

	return token, nil
}

func (s *ShareService) GetSharedObject(token string) (io.ReadCloser, *storage_models.Objects, error) {
	if token == "" {
		return nil, nil, errors.New("token is empty")
	}

	share, err := s.repo.GetObjectShareByToken(token)
	if err != nil || share == nil {
		return nil, nil, errors.New("share link not found or invalid")
	}

	if share.Revoked {
		return nil, nil, errors.New("this share link has been revoked")
	}
	if share.ExpiresAt != nil && share.ExpiresAt.Before(time.Now()) {
		return nil, nil, errors.New("this share link has expired")
	}

	if share.Object.ID == uuid.Nil || share.Object.Bucket.ID == uuid.Nil {
		return nil, nil, errors.New("associated object or bucket not found")
	}

	bucketPath := share.Object.OwnerID.String() + "/" + share.Object.Bucket.Name
	stream, err := s.storage.Get(bucketPath, share.Object.FileName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to retrieve file from storage: %w", err)
	}

	return stream, &share.Object, nil
}

func (s *ShareService) ListUserShares(userID string) ([]storage_models.ObjectShares, error) {
	ownerUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	shares, err := s.repo.ListObjectSharesByOwner(ownerUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve shares: %w", err)
	}
	return shares, nil
}

func (s *ShareService) RevokeShare(token string) error {
	return nil
}
