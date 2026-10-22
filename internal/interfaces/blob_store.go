package interfaces

import (
	"context"
	"errors"
	"io"
	"time"
)

// ErrPresignNotSupported is returned by BlobStore implementations that cannot
// produce presigned URLs (e.g. the local filesystem). Callers should fall back
// to the proxy upload/download path.
var ErrPresignNotSupported = errors.New("presigned URLs are not supported by this storage driver")

// BlobInfo is the metadata of a stored blob.
type BlobInfo struct {
	Key         string
	Size        int64
	ContentType string
	ETag        string
}

// PostTarget describes a presigned POST target: the endpoint plus the form
// fields the client must submit alongside the file. Unlike a presigned PUT, a
// POST carries a policy document, so the server can enforce conditions such as
// content-length-range (maximum upload size).
type PostTarget struct {
	URL    string
	Fields map[string]string
}

// BlobStore is the key-based storage abstraction shared by every backend
// (local filesystem, RustFS, or any S3-compatible service).
type BlobStore interface {
	// Driver returns the storage driver name ("local" or "s3").
	Driver() string
	Put(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error
	Get(ctx context.Context, key string) (io.ReadCloser, error)
	Stat(ctx context.Context, key string) (BlobInfo, error)
	Delete(ctx context.Context, key string) error
	PresignPutURL(ctx context.Context, key string, contentType string, ttl time.Duration) (string, error)
	// PresignPostURL returns a presigned POST form. When maxBytes > 0 the policy
	// enforces content-length-range so oversized uploads are rejected by the
	// storage server before any bytes are persisted.
	PresignPostURL(ctx context.Context, key string, contentType string, maxBytes int64, ttl time.Duration) (PostTarget, error)
	PresignGetURL(ctx context.Context, key string, contentType string, disposition string, filename string, ttl time.Duration) (string, error)
}

// ObjectKey builds the canonical blob key shared by every driver:
// {ownerID}/{bucketName}/{storageName}. It mirrors the legacy local layout so
// the migration to S3 is a straight 1:1 copy.
func ObjectKey(ownerID, bucketName, storageName string) string {
	return ownerID + "/" + bucketName + "/" + storageName
}
