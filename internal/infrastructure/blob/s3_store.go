package blob

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"objexrg/internal/interfaces"
)

// S3BlobStore talks to any S3-compatible service (RustFS by default) through
// the AWS SDK for Go v2.
type S3BlobStore struct {
	client  *s3.Client
	presign *s3.PresignClient
	bucket  string
}

type S3Options struct {
	Endpoint       string
	PublicEndpoint string
	Region         string
	AccessKey      string
	SecretKey      string
	Bucket         string
	UseSSL         bool
	ForcePathStyle bool
}

func NewS3BlobStore(opts S3Options) (*S3BlobStore, error) {
	if opts.Bucket == "" {
		return nil, errors.New("S3_BUCKET is required")
	}
	if opts.AccessKey == "" || opts.SecretKey == "" {
		return nil, errors.New("S3_ACCESS_KEY and S3_SECRET_KEY are required")
	}

	endpoint := normalizeEndpoint(opts.Endpoint, opts.UseSSL)
	if endpoint == "" {
		return nil, errors.New("S3_ENDPOINT is required")
	}

	region := opts.Region
	if region == "" {
		region = "us-east-1"
	}

	awsCfg := aws.Config{
		Region:      region,
		Credentials: credentials.NewStaticCredentialsProvider(opts.AccessKey, opts.SecretKey, ""),
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = opts.ForcePathStyle
	})

	// Presigned URLs are handed to the browser, so they must use an endpoint
	// the browser can reach (localhost in Docker Compose). Server-side
	// operations keep using the internal endpoint.
	presignEndpoint := normalizeEndpoint(opts.PublicEndpoint, opts.UseSSL)
	if presignEndpoint == "" {
		presignEndpoint = endpoint
	}
	presignClient := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(presignEndpoint)
		o.UsePathStyle = opts.ForcePathStyle
	})

	return &S3BlobStore{
		client:  client,
		presign: s3.NewPresignClient(presignClient),
		bucket:  opts.Bucket,
	}, nil
}

// EnsureBucket creates the backing bucket when it does not exist yet.
func (s *S3BlobStore) EnsureBucket(ctx context.Context) error {
	if _, err := s.client.HeadBucket(ctx, &s3.HeadBucketInput{Bucket: aws.String(s.bucket)}); err == nil {
		return nil
	}
	if _, err := s.client.CreateBucket(ctx, &s3.CreateBucketInput{Bucket: aws.String(s.bucket)}); err != nil {
		return fmt.Errorf("ensure bucket %q: %w", s.bucket, err)
	}
	return nil
}

func (s *S3BlobStore) Driver() string { return DriverS3 }

func (s *S3BlobStore) Put(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error {
	input := &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   reader,
	}
	if contentType != "" {
		input.ContentType = aws.String(contentType)
	}
	if size >= 0 {
		input.ContentLength = aws.Int64(size)
	}

	if _, err := s.client.PutObject(ctx, input); err != nil {
		return fmt.Errorf("put object: %w", err)
	}
	return nil
}

func (s *S3BlobStore) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	out, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("get object: %w", err)
	}
	return out.Body, nil
}

func (s *S3BlobStore) Stat(ctx context.Context, key string) (interfaces.BlobInfo, error) {
	out, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return interfaces.BlobInfo{}, fmt.Errorf("stat object: %w", err)
	}

	info := interfaces.BlobInfo{Key: key}
	if out.ContentLength != nil {
		info.Size = *out.ContentLength
	}
	if out.ContentType != nil {
		info.ContentType = *out.ContentType
	}
	if out.ETag != nil {
		info.ETag = strings.Trim(*out.ETag, "\"")
	}
	return info, nil
}

func (s *S3BlobStore) Delete(ctx context.Context, key string) error {
	if _, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}); err != nil {
		return fmt.Errorf("delete object: %w", err)
	}
	return nil
}

func (s *S3BlobStore) PresignPutURL(ctx context.Context, key string, contentType string, ttl time.Duration) (string, error) {
	input := &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}
	// The content type is part of the signature, so the client MUST send the
	// same Content-Type header it requested here.
	if contentType != "" {
		input.ContentType = aws.String(contentType)
	}

	req, err := s.presign.PresignPutObject(ctx, input, presignOptions(ttl))
	if err != nil {
		return "", fmt.Errorf("presign put: %w", err)
	}
	return req.URL, nil
}

// PresignPostURL builds a presigned POST form whose policy enforces a maximum
// object size (content-length-range) and, when set, an exact Content-Type. This
// lets the storage server reject oversized or wrong-typed uploads itself, before
// the bytes are persisted, instead of only detecting them at finalize time.
func (s *S3BlobStore) PresignPostURL(ctx context.Context, key string, contentType string, maxBytes int64, ttl time.Duration) (interfaces.PostTarget, error) {
	input := &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}
	if contentType != "" {
		input.ContentType = aws.String(contentType)
	}

	conditions := make([]any, 0, 2)
	if maxBytes > 0 {
		conditions = append(conditions, []any{"content-length-range", int64(0), maxBytes})
	}
	if contentType != "" {
		conditions = append(conditions, []any{"eq", "$Content-Type", contentType})
	}

	req, err := s.presign.PresignPostObject(ctx, input, func(o *s3.PresignPostOptions) {
		if ttl > 0 {
			o.Expires = ttl
		}
		if len(conditions) > 0 {
			o.Conditions = conditions
		}
	})
	if err != nil {
		return interfaces.PostTarget{}, fmt.Errorf("presign post: %w", err)
	}

	fields := req.Values
	if fields == nil {
		fields = map[string]string{}
	}
	// The Content-Type is a signed policy condition, so it must travel as a form
	// field exactly as requested here.
	if contentType != "" {
		fields["Content-Type"] = contentType
	}

	return interfaces.PostTarget{URL: req.URL, Fields: fields}, nil
}

func (s *S3BlobStore) PresignGetURL(ctx context.Context, key string, contentType string, disposition string, filename string, ttl time.Duration) (string, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}
	if contentType != "" {
		input.ResponseContentType = aws.String(contentType)
	}
	if disposition != "" {
		value := disposition
		if filename != "" {
			value = disposition + `; filename="` + filename + `"`
		}
		input.ResponseContentDisposition = aws.String(value)
	}

	req, err := s.presign.PresignGetObject(ctx, input, presignOptions(ttl))
	if err != nil {
		return "", fmt.Errorf("presign get: %w", err)
	}
	return req.URL, nil
}

func presignOptions(ttl time.Duration) func(*s3.PresignOptions) {
	return func(o *s3.PresignOptions) {
		if ttl > 0 {
			o.Expires = ttl
		}
	}
}

func normalizeEndpoint(endpoint string, useSSL bool) string {
	endpoint = strings.TrimSpace(endpoint)
	if endpoint == "" {
		return ""
	}
	if strings.HasPrefix(endpoint, "http://") || strings.HasPrefix(endpoint, "https://") {
		return endpoint
	}
	scheme := "http"
	if useSSL {
		scheme = "https"
	}
	return scheme + "://" + endpoint
}
