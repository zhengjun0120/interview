package response

import "github.com/gin-gonic/gin"

type JsonMsgResponse struct {
	Ctx *gin.Context
}

type JsonMsgResult struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type nilStruct struct {
}

const (
	SUCCESS_CODE = 200
	SUCCESS_MSG  = "成功"
)

func NewResponse(ctx *gin.Context) *JsonMsgResponse {
	return &JsonMsgResponse{Ctx: ctx}
}

func (r *JsonMsgResponse) Success(data interface{}) {
	res := JsonMsgResult{
		Code:    SUCCESS_CODE,
		Message: SUCCESS_MSG,
		Data:    data,
	}
	r.Ctx.JSON(200, res)
}

func (r *JsonMsgResponse) Error(mc MsgCode) {
	res := JsonMsgResult{
		Code:    mc.Code,
		Message: mc.Msg,
		Data:    nilStruct{},
	}
	r.Ctx.JSON(200, res)
}

func (r *JsonMsgResponse) ErrorWithStatus(mc MsgCode, status int) {
	res := JsonMsgResult{
		Code:    mc.Code,
		Message: mc.Msg,
		Data:    nilStruct{},
	}
	r.Ctx.JSON(status, res)
}
