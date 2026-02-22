package types

import (
	"ai_interview/biz/entity"
	"context"
	"time"
)

type IJobProfileService interface {
	SaveJobProfile(ctx context.Context, req *SaveJobProfileParams) error

	GetJobProfile(ctx context.Context) ([]GetJobProfile, error)
}

type SaveJobProfileParams struct {
	JobTitle               string
	Competencies           []entity.CompetenciesJson
	RedLineCondition       []string
	AiAdjustmentSuggestion string
}

type GetJobProfile struct {
	JobProfileID           string
	JobTitle               string
	Competencies           []entity.CompetenciesJson
	RedLineCondition       []string
	AiAdjustmentSuggestion string
	CreatedAt              time.Time
}
