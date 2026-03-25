package cos_service

import (
	"context"
	"mime/multipart"
)

type ICosService interface {
	UploadResume(file multipart.File, fileHeader *multipart.FileHeader, resumeID string) (string, error)
	DeleteFile(ctx context.Context, resumeUrl string)
}
