package response

import "time"

type BucketInfo struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

type BucketsListResponse struct {
	Buckets []BucketInfo `json:"buckets"`
}