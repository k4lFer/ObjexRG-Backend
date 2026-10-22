package storage

// NoopBucketStorage is the bucket backend used with S3-compatible drivers.
// In object storage every "bucket" is a key prefix, so it is created lazily on
// the first upload and needs no filesystem bookkeeping.
type NoopBucketStorage struct{}

func NewNoopBucketStorage() *NoopBucketStorage {
	return &NoopBucketStorage{}
}

func (n *NoopBucketStorage) CreateRootBucket(_ string) error { return nil }

func (n *NoopBucketStorage) CreateSubBucket(_, _ string) error { return nil }

func (n *NoopBucketStorage) DeleteSubBucket(_, _ string) error { return nil }

func (n *NoopBucketStorage) BucketExists(_, _ string) bool { return true }
