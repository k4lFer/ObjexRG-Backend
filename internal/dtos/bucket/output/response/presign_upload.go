package response

type PresignUploadResponse struct {
	Mode       string            `json:"mode"`
	Method     string            `json:"method,omitempty"`
	UploadURL  string            `json:"uploadUrl,omitempty"`
	StorageKey string            `json:"storageKey"`
	Key        string            `json:"key"`
	Headers    map[string]string `json:"headers,omitempty"`
	// Fields are the form fields to submit (plus the file) for a presigned POST.
	Fields map[string]string `json:"fields,omitempty"`
}
