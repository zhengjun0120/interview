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
		r.Error(response.MsgCode{Code: -1, Msg: "test"})
		return
	}

	params := caster.CastLoginReq2ServiceParams(&req)

	serviceResp, err := h.UserServer.Login(ctx, params)
	resp := caster.CastServiceResp2LoginResp(serviceResp)
	if err != nil {
		msgCode := ErrorToMsgCode(err)
		if msgCode == response.COMMON_FAIL {
			msgCode.Msg = err.Error()
		}
		r.Error(msgCode)
	} else {
		r.Success(resp)
	}
}
