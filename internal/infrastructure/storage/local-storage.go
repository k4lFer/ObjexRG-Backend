package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	basePath      string
	bucketStorage *BucketStorage
}

func NewLocalStorage(basePath string) *LocalStorage {
	return &LocalStorage{
		basePath:      basePath,
		bucketStorage: NewBucketStorage(basePath),
	}
}

func (l *LocalStorage) bucketPath(bucket string) string {
	return filepath.Join(l.basePath, bucket)
}

func (l *LocalStorage) rootBucketPath(userID string) string {
	return filepath.Join(l.basePath, userID)
}

func (l *LocalStorage) subBucketPath(userID, subBucketName string) string {
	return filepath.Join(l.basePath, userID, subBucketName)
}

func (l *LocalStorage) objectPath(bucket, filename string) string {
	return filepath.Join(l.bucketPath(bucket), filename)
}

func (l *LocalStorage) CreateRootBucket(userID string) error {
	return l.bucketStorage.CreateRootBucket(userID)
}

func (l *LocalStorage) CreateSubBucket(userID, subBucketName string) error {
	return l.bucketStorage.CreateSubBucket(userID, subBucketName)
}

func (l *LocalStorage) DeleteSubBucket(userID, subBucketName string) error {
	return l.bucketStorage.DeleteSubBucket(userID, subBucketName)
}

func (l *LocalStorage) BucketExists(userID, bucketName string) bool {
	return l.bucketStorage.BucketExists(userID, bucketName)
}

func (l *LocalStorage) Create(bucket string) error {
	path := l.bucketPath(bucket)
	return os.MkdirAll(path, 0755)
}

func (l *LocalStorage) Delete(bucket string) error {
	path := l.bucketPath(bucket)
	return os.RemoveAll(path)
}

func (l *LocalStorage) Save(bucket string, file io.Reader, filename *string) (string, error) {
	if bucket == "" {
		return "", fmt.Errorf("bucket is required")
	}
	if filename == nil || *filename == "" {
		return "", fmt.Errorf("filename is required")
	}

	// Verificar que el bucket existe
	bucketPath := l.bucketPath(bucket)
	if _, err := os.Stat(bucketPath); os.IsNotExist(err) {
		return "", fmt.Errorf("bucket does not exist")
	}

	target := l.objectPath(bucket, *filename)
	f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return "", fmt.Errorf("abrir archivo para escribir: %w", err)
	}
	defer f.Close()

	bytesWritten, err := io.Copy(f, file)
	if err != nil {
		return "", fmt.Errorf("escribir archivo %q: %w", target, err)
	}

	return fmt.Sprintf("%s (bytes:%d)", target, bytesWritten), nil
}

func (l *LocalStorage) Get(bucket string, filename string) (io.ReadCloser, error) {
	if bucket == "" || filename == "" {
		return nil, fmt.Errorf("bucket y filename son requeridos")
	}
	path := l.objectPath(bucket, filename)
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("abrir objeto: %w", err)
	}
	return file, nil
}

func (l *LocalStorage) ListBuckets() ([]string, error) {
	entries, err := os.ReadDir(l.basePath)
	if err != nil {
		return nil, fmt.Errorf("leer directorio base: %w", err)
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

func (l *LocalStorage) ListUserBuckets(userID string) ([]string, error) {
	return l.bucketStorage.ListUserBuckets(userID)
}

func (l *LocalStorage) ListObjects(bucket string) ([]string, error) {
	return l.bucketStorage.ListObjects(bucket)
}
