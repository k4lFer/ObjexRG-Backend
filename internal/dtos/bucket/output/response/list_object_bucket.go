package bucket_response

import "time"

type ObjectInfo struct {
	ID        string    `json:"id"`
	FileName  string    `json:"fileName"`
	SizeBytes int64     `json:"sizeBytes"`
	Extension string    `json:"extension"`
	MimeType  string    `json:"mimeType"`
	IsPublic  bool      `json:"isPublic"`
	CreatedAt time.Time `json:"createdAt"`
}

type ObjectsListResponse struct {
	Bucket  string       `json:"bucket"`
	Objects []ObjectInfo `json:"objects"`
}