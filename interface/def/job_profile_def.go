package def

import (
	"ai_interview/biz/entity"
	"time"
)

type SaveJobProfileRequest struct {
	JobTitle               string                    `json:"job_title"`
	Competencies           []entity.CompetenciesJson `json:"competencies"`
	RedLineCondition       []string                  `json:"red_line_condition"`
	AiAdjustmentSuggestion string                    `json:"ai_adjustment_suggestion"`
}

type GetJobProfile struct {
	JobProfileID           string                    `json:"job_profile_id"`
	JobTitle               string                    `json:"job_title"`
	Competencies           []entity.CompetenciesJson `json:"competencies"`
	RedLineCondition       []string                  `json:"red_line_condition"`
	AiAdjustmentSuggestion string                    `json:"ai_adjustment_suggestion"`
	CreatedAt              time.Time                 `json:"created_at"`
}

type GetJobProfileResponse struct {
	List []GetJobProfile `json:"list"`
}
