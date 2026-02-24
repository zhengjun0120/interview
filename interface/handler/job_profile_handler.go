package handler

import (
	"ai_interview/interface/caster"
	"ai_interview/interface/def"
	"ai_interview/pkg/response"
	"ai_interview/pkg/zlog"
	"github.com/gin-gonic/gin"
)

func (h *Handler) SaveJobProfile(gCtx *gin.Context) {
	ctx := gCtx.Request.Context()

	r := response.NewResponse(gCtx)

	var req def.SaveJobProfileRequest
	if err := gCtx.ShouldBindJSON(&req); err != nil {
		r.Error(response.PARAM_ERROR)
		zlog.Errorf("SaveJobProfile接口参数错误: %v", err)
		return
	}

	err := h.JobProfileServer.SaveJobProfile(ctx, caster.CastSaveJobProfileReq2ServiceParams(&req))
	if err != nil {
		msgCode := ErrorToMsgCode(err)
		if msgCode == response.COMMON_FAIL {
			msgCode.Msg = err.Error()
		}
		r.Error(msgCode)
		zlog.Errorf("SaveJobProfile接口调用失败: %v", err)
		return
	} else {
		r.Success(nil)
	}

}

func (h *Handler) GetJobProfile(gCtx *gin.Context) {
	ctx := gCtx.Request.Context()
	r := response.NewResponse(gCtx)

	serviceResp, err := h.JobProfileServer.GetJobProfile(ctx)
	if err != nil {
		msgCode := ErrorToMsgCode(err)
		r.Error(msgCode)
		zlog.Errorf("GetJobProfile接口调用失败: %v", err)
		return
	} else {
		r.Success(caster.CastServiceResp2GetJobProfileResp(serviceResp))
	}

}
