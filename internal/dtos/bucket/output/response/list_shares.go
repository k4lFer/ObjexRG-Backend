package response

import "time"

type ShareItemResponse struct {
	TokenHash     string     `json:"token"`
	Url           string     `json:"url"`
	BucketName    string     `json:"bucket_name"`
	FileName      string     `json:"file_name"`
	MimeType      string     `json:"mime_type"`
	SizeBytes     int64      `json:"size_bytes"`
	ExpiresAt     *time.Time `json:"expires_at"`
	Revoked       bool       `json:"revoked"`
	DownloadCount int64      `json:"download_count"`
	CreatedAt     time.Time  `json:"created_at"`
}

type ShareListResponse struct {
	Details []ShareItemResponse `json:"details"`
}
