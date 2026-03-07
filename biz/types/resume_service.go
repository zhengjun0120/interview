package types

import (
	"context"
	"mime/multipart"
	"time"
)

type IResumeService interface {

	// 上传简历
	UploadResume(ctx context.Context, req *UploadResumeParams) (*UploadResumeResponse, error)

	// 获取所有人才
	GetTalentAll(ctx context.Context) ([]GetTalentAllResponse, error)

	// 获取已面试人才
	GetTalentInterview(ctx context.Context) ([]GetTalentInterviewResponse, error)

	// 获取简历下载链接
	GetResumeUrl(ctx context.Context, req *GetResumeUrlRequest) (*GetResumeUrlResponse, error)

	// 获取人才报告
	GetTalentReport(ctx context.Context, req *GetTalentReportRequest)
}

type UploadResumeParams struct {
	File       multipart.File
	FileHeader *multipart.FileHeader
}

type UploadResumeResponse struct {
	TalentID        string
	FullName        string
	TargetPosition  string
	MatchScore      int
	InterviewStatus string
	CoreAdvantages  string
	HireStatus      string
	CreatedAt       time.Time
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

type GetResumeUrlRequest struct {
	TalentID string
}

type GetResumeUrlResponse struct {
	ResumeUrl string
}

type GetTalentReportRequest struct {
	TalentID string
}

type GetTalentReportQuestionArrJson struct {
	Question   string `json:"question"`
	Answer     string `json:"answer"`
	AiAnalysis string `json:"ai_analysis"`
}
type GetTalentReportResponse struct {
	List []GetTalentReportQuestionArrJson
}
