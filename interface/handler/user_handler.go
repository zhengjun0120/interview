package handler

import (
	"ai_interview/biz/types"
	"ai_interview/interface/caster"
	"ai_interview/interface/def"
	"ai_interview/pkg/response"
	"ai_interview/pkg/zlog"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Login(gCtx *gin.Context) {
	ctx := gCtx.Request.Context()
	var req def.LoginReq

	r := response.NewResponse(gCtx)

	if err := gCtx.ShouldBindJSON(&req); err != nil {
		r.Error(response.PARAM_ERROR)
		return
	}

	params := caster.CastLoginReq2ServiceParams(&req)

	serviceResp, err := h.UserServer.Login(ctx, params)

	if err != nil {
		msgCode := ErrorToMsgCode(err)
		if msgCode == response.COMMON_FAIL {
			msgCode.Msg = err.Error()
		}
		r.Error(msgCode)
	} else {
		resp := caster.CastServiceResp2LoginResp(serviceResp)
		r.Success(resp)
	}
}

func (h *Handler) Register(gCtx *gin.Context) {
	ctx := gCtx.Request.Context()
	var req def.RegisterReq

	r := response.NewResponse(gCtx)

	if err := gCtx.ShouldBindJSON(&req); err != nil {
		r.Error(response.PARAM_ERROR)
		return
	}

	params := caster.CastRegisterReq2ServiceParams(&req)
	serviceResp, err := h.UserServer.Register(ctx, params)

	if err != nil {
		msgCode := ErrorToMsgCode(err)
		if msgCode == response.COMMON_FAIL {
			msgCode.Msg = err.Error()
		}
		r.Error(msgCode)
	} else {
		resp := caster.CastServiceResp2RegisterResp(serviceResp)
		r.Success(resp)
	}
}

func (h *Handler) SendCode(gCtx *gin.Context) {
	ctx := gCtx.Request.Context()
	var req def.RegisterReq

	r := response.NewResponse(gCtx)

	if err := gCtx.ShouldBindJSON(&req); err != nil {
		r.Error(response.PARAM_ERROR)
		return
	}

	err := h.CodeServer.CaptchaSend(ctx, types.CaptchaWayTypeRegister, req.Email)
	if err != nil {
		msgCode := ErrorToMsgCode(err)
		if msgCode == response.COMMON_FAIL {
			msgCode.Msg = err.Error()
		}
		r.Error(msgCode)
		zlog.Errorf("SendCode接口调用失败: %v", err)
		return
	} else {
		r.Success(nil)
	}
}
