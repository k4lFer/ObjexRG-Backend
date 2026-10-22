package response

type ObjectVisibilityResponse struct {
	ID        string `json:"id"`
	Bucket    string `json:"bucket"`
	FileName  string `json:"fileName"`
	IsPublic  bool   `json:"isPublic"`
	PublicURL string `json:"publicUrl,omitempty"`
}
