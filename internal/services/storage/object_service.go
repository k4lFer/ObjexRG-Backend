package storage

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	storage_repo "objexrg/internal/infrastructure/repositories/storage"
	"objexrg/internal/interfaces"
	storage_models "objexrg/internal/models/storage"

	"github.com/google/uuid"
)

type ObjectService struct {
	store      interfaces.BlobStore
	repo       *storage_repo.StorageRepository
	presignTTL time.Duration
}

func NewObjectService(store interfaces.BlobStore, repo *storage_repo.StorageRepository, presignTTL time.Duration) *ObjectService {
	return &ObjectService{store: store, repo: repo, presignTTL: presignTTL}
}

// SetVisibility toggles whether an object can be read through the public route.
func (s *ObjectService) SetVisibility(userID, bucketName, filename string, isPublic bool) (*storage_models.Objects, error) {
	bucket, ownerUUID, err := loadBucket(s.repo, userID, bucketName)
	if err != nil {
		return nil, err
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

// DeleteObject removes the blob from storage and its database record.
func (s *ObjectService) DeleteObject(ctx context.Context, userID, bucketName, filename string) error {
	bucket, _, err := loadBucket(s.repo, userID, bucketName)
	if err != nil {
		return err
	}

	obj, err := s.repo.GetObjectByFileName(bucket.ID, filename)
	if err != nil {
		return err
	}
	if obj == nil {
		return fmt.Errorf("object not found")
	}

	key := interfaces.ObjectKey(obj.OwnerID.String(), bucketName, obj.StorageName())
	if err := s.store.Delete(ctx, key); err != nil {
		return fmt.Errorf("failed to remove file from storage: %w", err)
	}

	return s.repo.DeleteObject(obj.OwnerID, bucket.ID, filename)
}

// DownloadMeta resolves an object owned by the user (no blob streaming).
func (s *ObjectService) DownloadMeta(_ context.Context, userID, bucketName, filename string) (*storage_models.Objects, error) {
	bucket, _, err := loadBucket(s.repo, userID, bucketName)
	if err != nil {
		return nil, err
	}

	obj, err := s.repo.GetObjectByFileName(bucket.ID, filename)
	if err != nil {
		return nil, err
	}
	if obj == nil {
		return nil, errors.New("object not found")
	}
	return obj, nil
}

// GetPublic resolves an object that has been explicitly marked as public.
func (s *ObjectService) GetPublic(_ context.Context, objectID string) (*storage_models.Objects, error) {
	id, err := uuid.Parse(objectID)
	if err != nil {
		return nil, errors.New("invalid object id")
	}

	obj, err := s.repo.GetObjectByID(id)
	if err != nil {
		return nil, err
	}
	if obj == nil || !obj.IsPublic {
		return nil, errors.New("object not found")
	}
	return obj, nil
}

// LinkURL returns a URL the client can fetch directly: a presigned URL on S3
// drivers, or the authenticated proxy path on the local driver.
func (s *ObjectService) LinkURL(ctx context.Context, userID, bucketName, filename, disposition string) (string, error) {
	obj, err := s.DownloadMeta(ctx, userID, bucketName, filename)
	if err != nil {
		return "", err
	}

	key := interfaces.ObjectKey(obj.OwnerID.String(), bucketName, obj.StorageName())
	url2, err := s.store.PresignGetURL(ctx, key, obj.MimeType, disposition, obj.FileName, s.presignTTL)
	if err != nil {
		if errors.Is(err, interfaces.ErrPresignNotSupported) {
			return "/download/" + url.PathEscape(bucketName) + "/" + url.PathEscape(filename), nil
		}
		return "", err
	}
	return url2, nil
}
