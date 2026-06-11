package response

type UploadResponse struct {
	Url      string `json:"url"`
	Path     string `json:"path"`
	Filename string `json:"filename"`
}
