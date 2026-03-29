package bucket_response

type ObjectsListResponse struct {
	Bucket  string   `json:"bucket"`
	Objects []string `json:"objects"`
}