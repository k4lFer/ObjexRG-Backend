package blob

import (
	"context"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"time"

	"objexrg/internal/interfaces"
)

// LocalBlobStore stores blobs on the local filesystem. Keys are relative paths
// such as {ownerID}/{bucketName}/{storageName}, mirroring the S3 layout so a
// later migration is a straight copy.
type LocalBlobStore struct {
	basePath string
}

func NewLocalBlobStore(basePath string) *LocalBlobStore {
	return &LocalBlobStore{basePath: basePath}
}

func (l *LocalBlobStore) Driver() string { return DriverLocal }

// resolve turns a key into an absolute path, rejecting anything that escapes
// the base directory (path traversal).
func (l *LocalBlobStore) resolve(key string) (string, error) {
	clean := strings.TrimPrefix(filepath.ToSlash(filepath.Clean("/"+key)), "/")
	if clean == "" || clean == "." || strings.HasPrefix(clean, "..") {
		return "", fmt.Errorf("invalid storage key")
	}

	baseAbs, err := filepath.Abs(l.basePath)
	if err != nil {
		return "", err
	}
	fullAbs, err := filepath.Abs(filepath.Join(baseAbs, filepath.FromSlash(clean)))
	if err != nil {
		return "", err
	}
	if fullAbs != baseAbs && !strings.HasPrefix(fullAbs, baseAbs+string(os.PathSeparator)) {
		return "", fmt.Errorf("invalid storage key")
	}
	return fullAbs, nil
}

func (l *LocalBlobStore) Put(_ context.Context, key string, reader io.Reader, _ int64, _ string) error {
	full, err := l.resolve(key)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}

	f, err := os.OpenFile(full, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("open file for writing: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, reader); err != nil {
		return fmt.Errorf("write blob: %w", err)
	}
	return nil
}

func (l *LocalBlobStore) Get(_ context.Context, key string) (io.ReadCloser, error) {
	full, err := l.resolve(key)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(full)
	if err != nil {
		return nil, fmt.Errorf("open blob: %w", err)
	}
	return f, nil
}

func (l *LocalBlobStore) Stat(_ context.Context, key string) (interfaces.BlobInfo, error) {
	full, err := l.resolve(key)
	if err != nil {
		return interfaces.BlobInfo{}, err
	}
	info, err := os.Stat(full)
	if err != nil {
		return interfaces.BlobInfo{}, fmt.Errorf("stat blob: %w", err)
	}

	contentType := mime.TypeByExtension(filepath.Ext(full))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	return interfaces.BlobInfo{Key: key, Size: info.Size(), ContentType: contentType}, nil
}

func (l *LocalBlobStore) Delete(_ context.Context, key string) error {
	full, err := l.resolve(key)
	if err != nil {
		return err
	}
	if err := os.Remove(full); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("delete blob: %w", err)
	}
	return nil
}

func (l *LocalBlobStore) PresignPutURL(_ context.Context, _ string, _ string, _ time.Duration) (string, error) {
	return "", interfaces.ErrPresignNotSupported
}

func (l *LocalBlobStore) PresignPostURL(_ context.Context, _ string, _ string, _ int64, _ time.Duration) (interfaces.PostTarget, error) {
	return interfaces.PostTarget{}, interfaces.ErrPresignNotSupported
}

func (l *LocalBlobStore) PresignGetURL(_ context.Context, _ string, _ string, _ string, _ string, _ time.Duration) (string, error) {
	return "", interfaces.ErrPresignNotSupported
}
