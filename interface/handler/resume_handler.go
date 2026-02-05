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
