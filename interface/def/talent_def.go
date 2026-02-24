package def

import "time"

type GetTalentAll struct {
	TalentID        string    `json:"talent_id"`
	FullName        string    `json:"full_name"`
	TargetPosition  string    `json:"target_position"`
	MatchScore      int       `json:"match_score"`
	InterviewStatus string    `json:"interview_status"`
	CoreAdvantages  string    `json:"core_advantages"`
	HireStatus      string    `json:"hire_status"`
	CreatedAt       time.Time `json:"created_at"`
}

type GetTalentAllResponse struct {
	List []GetTalentAll `json:"list"`
}

type GetTalentInterview struct {
	TalentID       string    `json:"talent_id"`
	FullName       string    `json:"full_name"`
	TargetPosition string    `json:"target_position"`
	CoreAdvantages string    `json:"core_advantages"`
	HireStatus     string    `json:"hire_status"`
	CreatedAt      time.Time `json:"created_at"`
	InterviewTime  int       `json:"interview_time"`
}

type GetTalentInterviewResponse struct {
	List []GetTalentInterview `json:"list"`
}
