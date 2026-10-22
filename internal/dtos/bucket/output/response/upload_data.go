package response

type UploadDataResponse struct {
	Bucket   string `json:"bucket"`
	Filename string `json:"filename"`
	Path     string `json:"path"`
	Size     int64  `json:"size"`
}