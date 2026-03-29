package interfaces

type Bucket interface {
	CreateRootBucket(userID string) error
	CreateSubBucket(userID, subBucketName string) error
	DeleteSubBucket(userID, subBucketName string) error
	BucketExists(userID, bucketName string) bool
}
