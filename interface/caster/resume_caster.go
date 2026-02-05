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
