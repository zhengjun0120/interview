package caster

import (
	"ai_interview/biz/types"
	"ai_interview/interface/def"
)

func CastUploadResumeReq2ServiceParams(req *def.UploadResumeRequest) *types.UploadResumeParams {
	return &types.UploadResumeParams{
		File:       req.File,
		FileHeader: req.FileHeader,
	}
}

func CastServiceResp2UploadResumeResp(resp *types.UploadResumeResponse) *def.UploadResumeResponse {
	return &def.UploadResumeResponse{
		ResumeUrl: resp.ResumeUrl,
	}
}

func CastServiceResp2GetTalentAllResp(resp []types.GetTalentAllResponse) *def.GetTalentAllResponse {
	var res def.GetTalentAllResponse
	talents := make([]def.GetTalentAll, len(resp))
	for i, v := range resp {
		talents[i] = def.GetTalentAll{
			TalentID:        v.TalentID,
			FullName:        v.FullName,
			TargetPosition:  v.TargetPosition,
			MatchScore:      v.MatchScore,
			InterviewStatus: v.InterviewStatus,
			CoreAdvantages:  v.CoreAdvantages,
			HireStatus:      v.HireStatus,
			CreatedAt:       v.CreatedAt,
		}
	}

	res.List = talents
	return res
}

func CastServiceResp2GetTalentInterviewResp(resp []types.GetTalentInterviewResponse) *def.GetTalentInterviewResponse {
	var res def.GetTalentInterviewResponse

	talents := make([]def.GetTalentInterview, len(resp))
	for i, v := range resp {
		talents[i] = def.GetTalentInterview{
			TalentID:       v.TalentID,
			FullName:       v.FullName,
			TargetPosition: v.TargetPosition,
			InterviewTime:  v.InterviewTime,
			CoreAdvantages: v.CoreAdvantages,
			HireStatus:     v.HireStatus,
			CreatedAt:      v.CreatedAt,
		}
	}

	res.List = talents
	return res

}
