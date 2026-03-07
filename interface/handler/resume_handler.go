package handler

import (
	"ai_interview/interface/caster"
	"ai_interview/interface/def"
	"ai_interview/pkg/response"
	"ai_interview/pkg/zlog"
	"github.com/gin-gonic/gin"
)

func (h *Handler) UploadResume(gCtx *gin.Context) {
	ctx := gCtx.Request.Context()
	var req def.UploadResumeRequest

	r := response.NewResponse(gCtx)

	file, fileHeader, err := gCtx.Request.FormFile("resume")
	if err != nil {
		r.Error(response.PARAM_ERROR)
		zlog.Errorf("UploadResume接口参数错误，%v", err)
		return
	}
	req.File = file
	req.FileHeader = fileHeader

	params := caster.CastUploadResumeReq2ServiceParams(&req)
	serviceResp, err := h.ResumeServer.UploadResume(ctx, params)

	if err != nil {
		msgCode := ErrorToMsgCode(err)
		if msgCode == response.COMMON_FAIL {
			msgCode.Msg = err.Error()
		}
		r.Error(msgCode)
		zlog.Errorf("UploadResume接口调用失败，%v", err)
	} else {
		resp := caster.CastServiceResp2UploadResumeResp(serviceResp)
		r.Success(resp)
	}
}

func (h *Handler) GetTalent(gCtx *gin.Context) {
	ctx := gCtx.Request.Context()

	r := response.NewResponse(gCtx)

	query := gCtx.Query("type")
	if query != "all" && query != "interview" {
		r.Error(response.QUERY_PARAM_ERROR)
		zlog.Errorf("GetTalent接口参数错误，type:%s", query)
		return
	}

	if query == "all" {
		serviceResp, err := h.ResumeServer.GetTalentAll(ctx)
		if err == nil {
			resp := caster.CastServiceResp2GetTalentAllResp(serviceResp)
			r.Success(resp)
			return
		} else {
			msgCode := ErrorToMsgCode(err)
			if msgCode == response.COMMON_FAIL {
				msgCode.Msg = err.Error()
			}
			zlog.Errorf("GetTalent接口调用失败，%v", err)
			r.Error(msgCode)
		}
	} else {
		serviceResp, err := h.ResumeServer.GetTalentInterview(ctx)
		if err == nil {
			resp := caster.CastServiceResp2GetTalentInterviewResp(serviceResp)
			r.Success(resp)
			return
		} else {
			msgCode := ErrorToMsgCode(err)
			if msgCode == response.COMMON_FAIL {
				msgCode.Msg = err.Error()
			}
			zlog.Errorf("GetTalent接口调用失败，%v", err)
			r.Error(msgCode)
		}
	}

}

func (h *Handler) GetResumeUrl(gCtx *gin.Context) {
	ctx := gCtx.Request.Context()
	r := response.NewResponse(gCtx)
	var req def.GetResumeUrlRequest
	req.TalentID = gCtx.Query("talent_id")
	if req.TalentID == "" {
		r.Error(response.PARAM_ERROR)
		return
	}

	serviceResp, err := h.ResumeServer.GetResumeUrl(ctx, caster.CastGetResumeUrlReq2ServiceParams(&req))
	if err != nil {
		msgCode := ErrorToMsgCode(err)
		if msgCode == response.COMMON_FAIL {
			msgCode.Msg = err.Error()
		}
		r.Error(msgCode)
		zlog.Errorf("GetResumeUrl接口调用失败，%v", err)
	} else {
		resp := caster.CastGetResumeUrlResp2ServiceResp(serviceResp)
		r.Success(resp)
	}
}
