package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	storage_repo "objexrg/internal/infrastructure/repositories/storage"
	"objexrg/internal/interfaces"
	storage_models "objexrg/internal/models/storage"

	"github.com/google/uuid"
)

// ErrObjectExists is returned when a file with the same name already exists
// in the target bucket.
var ErrObjectExists = errors.New("an object with that name already exists in this bucket")

// ErrUploadTooLarge is returned when a (declared or stored) blob exceeds the
// configured maximum upload size.
var ErrUploadTooLarge = errors.New("file exceeds the maximum allowed size")

// PresignResult describes how a client should upload a blob.
//
// Mode is one of:
//   - "s3-post": direct presigned POST (form fields) — enforces the size limit
//     at the storage server via content-length-range.
//   - "s3":      direct presigned PUT (headers) — used when POST is unavailable.
//   - "proxy":   no presigning; fall back to POST /upload/{bucket}.
type PresignResult struct {
	Mode       string
	UploadURL  string
	StorageKey string
	Key        string
	Headers    map[string]string
	Fields     map[string]string
}

type UploadService struct {
	store          interfaces.BlobStore
	repo           *storage_repo.StorageRepository
	maxUploadBytes int64
	presignTTL     time.Duration
}

func NewUploadService(store interfaces.BlobStore, repo *storage_repo.StorageRepository, maxUploadBytes int64, presignTTL time.Duration) *UploadService {
	return &UploadService{store: store, repo: repo, maxUploadBytes: maxUploadBytes, presignTTL: presignTTL}
}

// BucketExists reports whether the user owns a bucket with that name.
func (s *UploadService) BucketExists(userID, bucketName string) bool {
	bucket, _, err := loadBucket(s.repo, userID, bucketName)
	return err == nil && bucket != nil
}

// Upload stores a blob through the backend (proxy path). It is used by the
// classic multipart endpoint and always available, even without presigning.
func (s *UploadService) Upload(ctx context.Context, userID, bucketName, filename string, size int64, mimeType string, reader io.Reader) (string, error) {
	if filename == "" {
		return "", errors.New("filename is required")
	}

	bucket, ownerUUID, err := loadBucket(s.repo, userID, bucketName)
	if err != nil {
		return "", err
	}

	existing, err := s.repo.GetObjectByFileName(bucket.ID, filename)
	if err != nil {
		return "", err
	}
	if existing != nil {
		return "", ErrObjectExists
	}

	// Store the blob under a UUID-based key so the original name is only
	// metadata. This keeps the storage layer collision-free and opens the
	// door for content-addressed deduplication later.
	objectID := uuid.New()
	extension := filepath.Ext(filename)
	storageKey := objectID.String() + extension
	key := interfaces.ObjectKey(ownerUUID.String(), bucketName, storageKey)

	if err := s.store.Put(ctx, key, reader, size, mimeType); err != nil {
		return "", err
	}

	obj := &storage_models.Objects{
		ID:         objectID,
		BucketID:   bucket.ID,
		OwnerID:    ownerUUID,
		Path:       key,
		FileName:   filename,
		StorageKey: storageKey,
		SizeBytes:  size,
		MimeType:   mimeType,
		Extension:  extension,
		CreatedAt:  time.Now(),
	}

	if err := s.repo.CreateObject(obj); err != nil {
		// Roll back the blob so we don't leave orphaned files behind.
		_ = s.store.Delete(ctx, key)
		return "", fmt.Errorf("file saved but db record failed: %w", err)
	}

	return key, nil
}

// PresignUpload reserves a storage key and, when the driver supports it,
// returns a presigned URL the client can PUT to directly.
func (s *UploadService) PresignUpload(ctx context.Context, userID, bucketName, filename string, size int64, mimeType string) (*PresignResult, error) {
	if filename == "" {
		return nil, errors.New("filename is required")
	}
	if s.maxUploadBytes > 0 && size > s.maxUploadBytes {
		return nil, ErrUploadTooLarge
	}

	bucket, ownerUUID, err := loadBucket(s.repo, userID, bucketName)
	if err != nil {
		return nil, err
	}

	existing, err := s.repo.GetObjectByFileName(bucket.ID, filename)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrObjectExists
	}

	objectID := uuid.New()
	extension := filepath.Ext(filename)
	storageKey := objectID.String() + extension
	key := interfaces.ObjectKey(ownerUUID.String(), bucketName, storageKey)

	// Prefer a presigned POST: its policy enforces the maximum size
	// (content-length-range) at the storage server, so oversized uploads are
	// rejected before any bytes are persisted.
	if post, err := s.store.PresignPostURL(ctx, key, mimeType, s.maxUploadBytes, s.presignTTL); err == nil {
		return &PresignResult{
			Mode:       "s3-post",
			UploadURL:  post.URL,
			Fields:     post.Fields,
			StorageKey: storageKey,
			Key:        key,
		}, nil
	} else if !errors.Is(err, interfaces.ErrPresignNotSupported) {
		return nil, err
	}

	// Fall back to a presigned PUT when the driver supports signing but not POST.
	url, err := s.store.PresignPutURL(ctx, key, mimeType, s.presignTTL)
	if err != nil {
		if errors.Is(err, interfaces.ErrPresignNotSupported) {
			return &PresignResult{Mode: "proxy", StorageKey: storageKey, Key: key}, nil
		}
		return nil, err
	}

	result := &PresignResult{Mode: "s3", UploadURL: url, StorageKey: storageKey, Key: key}
	if mimeType != "" {
		result.Headers = map[string]string{"Content-Type": mimeType}
	}
	return result, nil
}

// FinalizeUpload records the metadata of a blob that was uploaded directly to
// storage, validating its location and size.
func (s *UploadService) FinalizeUpload(ctx context.Context, userID, bucketName, filename, storageKey string) (*storage_models.Objects, error) {
	if filename == "" || storageKey == "" {
		return nil, errors.New("fileName and storageKey are required")
	}

	bucket, ownerUUID, err := loadBucket(s.repo, userID, bucketName)
	if err != nil {
		return nil, err
	}

	if strings.ContainsAny(storageKey, "/\\") || strings.Contains(storageKey, "..") {
		return nil, errors.New("invalid storage key")
	}

	key := interfaces.ObjectKey(ownerUUID.String(), bucketName, storageKey)

	info, err := s.store.Stat(ctx, key)
	if err != nil {
		return nil, errors.New("uploaded file not found in storage")
	}
	if s.maxUploadBytes > 0 && info.Size > s.maxUploadBytes {
		_ = s.store.Delete(ctx, key)
		return nil, ErrUploadTooLarge
	}

	existing, err := s.repo.GetObjectByFileName(bucket.ID, filename)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		_ = s.store.Delete(ctx, key)
		return nil, ErrObjectExists
	}

	mimeType := info.ContentType
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	obj := &storage_models.Objects{
		ID:         uuid.New(),
		BucketID:   bucket.ID,
		OwnerID:    ownerUUID,
		Path:       key,
		FileName:   filename,
		StorageKey: storageKey,
		SizeBytes:  info.Size,
		MimeType:   mimeType,
		Extension:  filepath.Ext(filename),
		CreatedAt:  time.Now(),
	}

	if err := s.repo.CreateObject(obj); err != nil {
		_ = s.store.Delete(ctx, key)
		return nil, fmt.Errorf("failed to save object metadata: %w", err)
	}

	return obj, nil
}
