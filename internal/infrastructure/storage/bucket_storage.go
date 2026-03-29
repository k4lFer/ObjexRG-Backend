package storage

import (
	"fmt"
	"os"
	"path/filepath"
)

type BucketStorage struct {
	basePath string
}

func NewBucketStorage(basePath string) *BucketStorage {
	return &BucketStorage{basePath: basePath}
}

func (b *BucketStorage) rootBucketPath(userID string) string {
	return filepath.Join(b.basePath, userID)
}

func (b *BucketStorage) subBucketPath(userID, subBucketName string) string {
	return filepath.Join(b.basePath, userID, subBucketName)
}

func (b *BucketStorage) bucketPath(bucket string) string {
	return filepath.Join(b.basePath, bucket)
}

func (b *BucketStorage) CreateRootBucket(userID string) error {
	path := b.rootBucketPath(userID)
	return os.MkdirAll(path, 0755)
}

func (b *BucketStorage) CreateSubBucket(userID, subBucketName string) error {
	path := b.subBucketPath(userID, subBucketName)
	return os.MkdirAll(path, 0755)
}

func (b *BucketStorage) DeleteSubBucket(userID, subBucketName string) error {
	path := b.subBucketPath(userID, subBucketName)
	return os.RemoveAll(path)
}

func (b *BucketStorage) BucketExists(userID, bucketName string) bool {
	var path string
	if bucketName == "" {
		path = b.rootBucketPath(userID)
	} else {
		path = b.subBucketPath(userID, bucketName)
	}
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func (b *BucketStorage) ListUserBuckets(userID string) ([]string, error) {
	userPath := b.rootBucketPath(userID)
	entries, err := os.ReadDir(userPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("leer buckets de usuario: %w", err)
	}
	var buckets []string
	for _, entry := range entries {
		if entry.IsDir() {
			buckets = append(buckets, entry.Name())
		}
	}
	if buckets == nil {
		buckets = []string{}
	}
	return buckets, nil
}

func (b *BucketStorage) ListObjects(bucket string) ([]string, error) {
	if bucket == "" {
		return nil, fmt.Errorf("bucket es requerido")
	}
	path := b.bucketPath(bucket)
	entries, err := os.ReadDir(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("leer bucket: %w", err)
	}
	var objects []string
	for _, entry := range entries {
		if !entry.IsDir() {
			objects = append(objects, entry.Name())
		}
	}
	if objects == nil {
		objects = []string{}
	}
	return objects, nil
}
