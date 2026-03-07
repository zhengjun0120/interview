package handler

import (
	"ai_interview/interface/caster"
	"ai_interview/interface/def"
	"ai_interview/pkg/response"
	"ai_interview/pkg/zlog"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

func (h *Handler) CreateInterview(gCtx *gin.Context) {
	ctx := gCtx.Request.Context()
	r := response.NewResponse(gCtx)
	var req def.CreateInterviewRequest
	if err := gCtx.ShouldBindJSON(&req); err != nil {
		r.Error(response.PARAM_ERROR)
		return
	}

	serviceResp, err := h.ChatServer.CreateInterviewRoom(ctx, caster.CastCreateInterviewReq2ServiceParams(&req))
	if err != nil {
		msgCode := ErrorToMsgCode(err)
		if msgCode == response.COMMON_FAIL {
			msgCode.Msg = err.Error()
		}
		r.Error(msgCode)
		zlog.Errorf("CreateInterviewRoom接口调用失败，%v", err)
		return
	} else {
		r.Success(caster.CastServiceResp2CreateInterviewResp(serviceResp))
	}
}

func (h *Handler) CheckRoomPermission(gCtx *gin.Context) {
	ctx := gCtx.Request.Context()
	r := response.NewResponse(gCtx)
	var req def.CheckRoomPermissionRequest
	if err := gCtx.ShouldBindJSON(&req); err != nil {
		r.Error(response.PARAM_ERROR)
		return
	}
	serviceResp, err := h.ChatServer.CheckRoomPermission(ctx, caster.CastCheckRoomPermissionReq2ServiceParams(&req))
	if err != nil {
		msgCode := ErrorToMsgCode(err)
		if msgCode == response.COMMON_FAIL {
			msgCode.Msg = err.Error()
		}
		r.Error(msgCode)
		zlog.Errorf("CheckRoomPermission接口调用失败，%v", err)
		return
	} else {
		r.Success(caster.CastServiceResp2CheckRoomPermissionResp(serviceResp))
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) JoinInterviewRoom(gCtx *gin.Context) {
	ctx := gCtx.Request.Context()
	var req def.JoinInterviewRoomRequest

	req.RoomID = gCtx.Query("room_id")

	conn, err := upgrader.Upgrade(gCtx.Writer, gCtx.Request, nil)
	if err != nil {
		zlog.Errorf("websocket升级失败: %v", err)
		return
	}
	req.Conn = conn

	err = h.ChatServer.JoinInterviewRoom(ctx, caster.CastJoinInterviewRoomReq2ServiceParams(&req))
	if err != nil {
		closeMsg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, err.Error())
		_ = conn.WriteMessage(websocket.CloseMessage, closeMsg)
		_ = conn.Close()

		zlog.Errorf("JoinInterviewRoom接口调用失败，%v", err)
		return
	}

}

func (h *Handler) AddTagToMessage(gCtx *gin.Context) {
	ctx := gCtx.Request.Context()
	var req def.AddTagToMessageRequest
	r := response.NewResponse(gCtx)
	if err := gCtx.ShouldBindJSON(&req); err != nil {
		r.Error(response.PARAM_ERROR)
		return
	}

	err := h.ChatServer.AddTagToMessage(ctx, caster.CastAddTagToMessageReq2ServiceParams(&req))
	if err != nil {
		msgCode := ErrorToMsgCode(err)
		if msgCode == response.COMMON_FAIL {
			msgCode.Msg = err.Error()
		}
		r.Error(msgCode)
		zlog.Errorf("AddTagToMessage接口调用失败，%v", err)
		return
	} else {
		r.Success(nil)
	}
}

func (h *Handler) GetAiSuggestion(gCtx *gin.Context) {
	ctx := gCtx.Request.Context()
	var req def.GetAiSuggestionRequest
	r := response.NewResponse(gCtx)

	if err := gCtx.ShouldBindJSON(&req); err != nil {
		r.Error(response.PARAM_ERROR)
		return
	}

	serviceResp, err := h.ChatServer.GetAiSuggestion(ctx, caster.CastGetAiSuggestionReq2ServiceParams(&req))
	if err != nil {
		msgCode := ErrorToMsgCode(err)
		if msgCode == response.COMMON_FAIL {
			msgCode.Msg = err.Error()
		}
		r.Error(msgCode)
		zlog.Errorf("GetAiSuggestion接口调用失败，%v", err)
		return
	} else {
		r.Success(caster.CastServiceResp2GetAiSuggestionResp(serviceResp))
	}
}

// 结束面试
func (h *Handler) EndInterview(gCtx *gin.Context) {
	ctx := gCtx.Request.Context()
	r := response.NewResponse(gCtx)
	var req def.EndInterviewRequest
	if err := gCtx.ShouldBindJSON(&req); err != nil {
		r.Error(response.PARAM_ERROR)
		return
	}

	err := h.ChatServer.EndInterview(ctx, caster.CastEndInterviewReq2ServiceParams(&req))
	if err != nil {
		msgCode := ErrorToMsgCode(err)
		if msgCode == response.COMMON_FAIL {
			msgCode.Msg = err.Error()
		}
		r.Error(msgCode)
		zlog.Errorf("EndInterview接口调用失败，%v", err)
		return
	} else {
		r.Success(nil)
	}
}
