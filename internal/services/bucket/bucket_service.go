package bucket

import (
	"fmt"
	"objexrg-backend/internal/interfaces"

	"gorm.io/gorm"
)

type BucketService struct {
	bucket interfaces.BucketInterface
	db	 *gorm.DB
}

func NewBucketService(bucket interfaces.BucketInterface, db *gorm.DB) *BucketService {
	return &BucketService{bucket: bucket, db: db}
}

func (s *BucketService) CreateRootBucket(userID string) error {
	return s.bucket.CreateRootBucket(userID)
}

func (s *BucketService) CreateSubBucket(userID, subBucketName string) error {
	if s.bucket.BucketExists(userID, subBucketName) {
		return fmt.Errorf("bucket %s already exists for user %s", subBucketName, userID)
	}
	return s.bucket.CreateSubBucket(userID, subBucketName)
}

func (s *BucketService) DeleteSubBucket(userID, subBucketName string) error {
	return s.bucket.DeleteSubBucket(userID, subBucketName)
}
