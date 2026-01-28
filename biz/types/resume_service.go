package types

import (
	"context"
	"mime/multipart"
)

type IResumeService interface {
	UploadResume(ctx context.Context, req *UploadResumeParams) (*UploadResumeResponse, error)
}

type UploadResumeParams struct {
	File       multipart.File
	FileHeader *multipart.FileHeader
}

type UploadResumeResponse struct {
	ResumeUrl string
}
