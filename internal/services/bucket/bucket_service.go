package bucket

import (
	"fmt"
	"time"

	bucket_response "objexrg-backend/internal/dtos/bucket/output/response"
	storage_repo "objexrg-backend/internal/infrastructure/repositories/storage"
	"objexrg-backend/internal/interfaces"
	storage_models "objexrg-backend/internal/models/storage"

	"github.com/google/uuid"
)

type BucketService struct {
	bucket interfaces.BucketInterface
	repo   *storage_repo.StorageRepository
}

func NewBucketService(bucket interfaces.BucketInterface, repo *storage_repo.StorageRepository) *BucketService {
	return &BucketService{bucket: bucket, repo: repo}
}

func (s *BucketService) CreateRootBucket(userID string) error {
	return s.bucket.CreateRootBucket(userID)
}

func (s *BucketService) CreateSubBucket(userID, subBucketName string) error {
	if s.bucket.BucketExists(userID, subBucketName) {
		return fmt.Errorf("bucket %s already exists for user %s", subBucketName, userID)
	}

	ownerUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	// Create DB Record first
	dbBucket := &storage_models.Bucket{
		OwnerID:   ownerUUID,
		Name:      subBucketName,
		Status:    "active",
		CreatedAt: time.Now(),
	}

	if err := s.repo.CreateBucket(dbBucket); err != nil {
		return fmt.Errorf("failed to save bucket in db: %w", err)
	}

	// Create FileSystem Directory
	if err := s.bucket.CreateSubBucket(userID, subBucketName); err != nil {
		// Rollback could be implemented here
		s.repo.DeleteBucketByOwner(ownerUUID, subBucketName)
		return err
	}

	return nil
}

func (s *BucketService) DeleteSubBucket(userID, subBucketName string) error {
	ownerUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	if err := s.repo.DeleteBucketByOwner(ownerUUID, subBucketName); err != nil {
		return fmt.Errorf("failed to delete bucket in db: %w", err)
	}

	return s.bucket.DeleteSubBucket(userID, subBucketName)
}

func (s *BucketService) ListBuckets(userID string) ([]bucket_response.BucketInfo, error) {
	ownerUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	dbBuckets, err := s.repo.ListBucketsByOwner(ownerUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to list buckets from db: %w", err)
	}

	var infos []bucket_response.BucketInfo
	for _, b := range dbBuckets {
		infos = append(infos, bucket_response.BucketInfo{
			ID:        b.ID.String(),
			Name:      b.Name,
			Status:    b.Status,
			CreatedAt: b.CreatedAt,
		})
	}
	return infos, nil
}

func (s *BucketService) ListObjects(userID, bucketName string) ([]bucket_response.ObjectInfo, error) {
	ownerUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	bucket, err := s.repo.GetBucketByNameAndOwner(ownerUUID, bucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to query bucket: %w", err)
	}
	if bucket == nil {
		return nil, fmt.Errorf("bucket not found")
	}

	dbObjects, err := s.repo.ListObjectsByBucket(bucket.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to list objects from db: %w", err)
	}

	var infos []bucket_response.ObjectInfo
	for _, obj := range dbObjects {
		infos = append(infos, bucket_response.ObjectInfo{
			ID:        obj.ID.String(),
			FileName:  obj.FileName,
			SizeBytes: obj.SizeBytes,
			Extension: obj.Extension,
			MimeType:  obj.MimeType,
			CreatedAt: obj.CreatedAt,
		})
	}
	return infos, nil
}

