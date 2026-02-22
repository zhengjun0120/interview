package types

import (
	"context"
	"mime/multipart"
	"time"
)

type IResumeService interface {
	UploadResume(ctx context.Context, req *UploadResumeParams) (*UploadResumeResponse, error)

	GetTalentAll(ctx context.Context) ([]GetTalentAllResponse, error)

	GetTalentInterview(ctx context.Context) ([]GetTalentInterviewResponse, error)
}

type UploadResumeParams struct {
	File       multipart.File
	FileHeader *multipart.FileHeader
}

type UploadResumeResponse struct {
	ResumeUrl string
}

type GetTalentAllResponse struct {
	TalentID        string
	FullName        string
	TargetPosition  string
	MatchScore      int
	InterviewStatus string
	CoreAdvantages  string
	HireStatus      string
	CreatedAt       time.Time
}

type GetTalentInterviewResponse struct {
	TalentID       string
	FullName       string
	TargetPosition string
	CoreAdvantages string
	HireStatus     string
	CreatedAt      time.Time
	InterviewTime  int
}
