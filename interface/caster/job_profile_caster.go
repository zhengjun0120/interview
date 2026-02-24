package caster

import (
	"ai_interview/biz/types"
	"ai_interview/interface/def"
)

func CastSaveJobProfileReq2ServiceParams(req *def.SaveJobProfileRequest) *types.SaveJobProfileParams {
	return &types.SaveJobProfileParams{
		JobTitle:               req.JobTitle,
		RedLineCondition:       req.RedLineCondition,
		AiAdjustmentSuggestion: req.AiAdjustmentSuggestion,
		Competencies:           req.Competencies,
	}
}

func CastServiceResp2GetJobProfileResp(resp []types.GetJobProfile) *def.GetJobProfileResponse {

	res := make([]def.GetJobProfile, len(resp))

	for i, v := range resp {
		res[i] = def.GetJobProfile{
			JobProfileID:           v.JobProfileID,
			JobTitle:               v.JobTitle,
			RedLineCondition:       v.RedLineCondition,
			AiAdjustmentSuggestion: v.AiAdjustmentSuggestion,
			Competencies:           v.Competencies,
			CreatedAt:              v.CreatedAt,
		}
	}
	return &def.GetJobProfileResponse{List: res}
}
