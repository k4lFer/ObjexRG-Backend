package response

import "time"

type ObjectUploadResponse struct {
	ID         string    `json:"id"`
	Bucket     string    `json:"bucket"`
	FileName   string    `json:"fileName"`
	SizeBytes  int64     `json:"sizeBytes"`
	MimeType   string    `json:"mimeType"`
	StorageKey string    `json:"storageKey"`
	Key        string    `json:"key"`
	CreatedAt  time.Time `json:"createdAt"`
}
