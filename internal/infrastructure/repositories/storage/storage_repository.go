package storage_repository

import (
	"errors"
	storage_models "objexrg-backend/internal/models/storage"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StorageRepository struct {
	db *gorm.DB
}

func NewStorageRepository(db *gorm.DB) *StorageRepository {
	return &StorageRepository{db: db}
}

func (r *StorageRepository) CreateBucket(bucket *storage_models.Bucket) error {
	return r.db.Create(bucket).Error
}

func (r *StorageRepository) GetBucketByNameAndOwner(ownerID uuid.UUID, name string) (*storage_models.Bucket, error) {
	var bucket storage_models.Bucket
	err := r.db.Where("owner_id = ? AND name = ?", ownerID, name).First(&bucket).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Bucket doesn't exist
		}
		return nil, err
	}
	return &bucket, nil
}

func (r *StorageRepository) DeleteBucketByOwner(ownerID uuid.UUID, name string) error {
	return r.db.Where("owner_id = ? AND name = ?", ownerID, name).Delete(&storage_models.Bucket{}).Error
}

func (r *StorageRepository) ListBucketsByOwner(ownerID uuid.UUID) ([]storage_models.Bucket, error) {
	var buckets []storage_models.Bucket
	err := r.db.Where("owner_id = ?", ownerID).Find(&buckets).Error
	return buckets, err
}

func (r *StorageRepository) CreateObject(obj *storage_models.Objects) error {
	return r.db.Create(obj).Error
}

func (r *StorageRepository) DeleteObject(ownerID uuid.UUID, bucketID uuid.UUID, fileName string) error {
	return r.db.Where("owner_id = ? AND bucket_id = ? AND file_name = ?", ownerID, bucketID, fileName).Delete(&storage_models.Objects{}).Error
}

func (r *StorageRepository) ListObjectsByBucket(bucketID uuid.UUID) ([]storage_models.Objects, error) {
	var objects []storage_models.Objects
	err := r.db.Where("bucket_id = ?", bucketID).Find(&objects).Error
	return objects, err
}

func (r *StorageRepository) GetObjectByFileName(bucketID uuid.UUID, fileName string) (*storage_models.Objects, error) {
	var obj storage_models.Objects
	err := r.db.Where("bucket_id = ? AND file_name = ?", bucketID, fileName).First(&obj).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &obj, nil
}

func (r *StorageRepository) CreateObjectShare(share *storage_models.ObjectShares) error {
	return r.db.Create(share).Error
}

func (r *StorageRepository) GetObjectShareByToken(token string) (*storage_models.ObjectShares, error) {
	var share storage_models.ObjectShares
	err := r.db.Preload("Object").Preload("Object.Bucket").Where("token_hash = ?", token).First(&share).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &share, nil
}

func (r *StorageRepository) ListObjectSharesByOwner(ownerID uuid.UUID) ([]storage_models.ObjectShares, error) {
	var shares []storage_models.ObjectShares
	err := r.db.Preload("Object").Preload("Object.Bucket").Where("owner_id = ?", ownerID).Find(&shares).Error
	return shares, err
}
