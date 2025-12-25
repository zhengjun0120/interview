package handler

import (
	"ai_interview/interface/caster"
	"ai_interview/interface/def"
	"ai_interview/pkg/response"
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
