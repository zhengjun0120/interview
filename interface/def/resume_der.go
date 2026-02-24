package def

import (
	"mime/multipart"
	"time"
)

type UploadResumeRequest struct {
	File       multipart.File
	FileHeader *multipart.FileHeader
}

type UploadResumeResponse struct {
	TalentID        string    `json:"talent_id"`
	FullName        string    `json:"full_name"`
	TargetPosition  string    `json:"target_position"`
	MatchScore      int       `json:"match_score"`
	InterviewStatus string    `json:"interview_status"`
	CoreAdvantages  string    `json:"core_advantages"`
	HireStatus      string    `json:"hire_status"`
	CreatedAt       time.Time `json:"created_at"`
}
