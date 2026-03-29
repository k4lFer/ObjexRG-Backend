package bucket_response

type CreateBucketResponse struct {
	Bucket string `json:"bucket"`
	User   string `json:"user"`
}