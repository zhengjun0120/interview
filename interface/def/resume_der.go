package def

import "mime/multipart"

type UploadResumeRequest struct {
	File       multipart.File
	FileHeader *multipart.FileHeader
}

type UploadResumeResponse struct {
	ResumeUrl string `json:"resume_url"`
}
