package interfaces

import "io"

type StorageInterface interface {
	Save(bucket string, file io.Reader, filename *string) (string, error)
	Get(bucket string, filename string) (io.ReadCloser, error)
	DeleteObject(bucket string, filename string) error
	ListBuckets() ([]string, error)
	ListUserBuckets(userID string) ([]string, error)
	ListObjects(bucket string) ([]string, error)
	BucketExists(userID, bucketName string) bool
}
