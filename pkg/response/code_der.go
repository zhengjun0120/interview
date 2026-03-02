package response

type MsgCode struct {
	Code int
	Msg  string
}

var (
	//成功响应
	SUCCESS = MsgCode{Code: 200, Msg: "成功"}

	//通用错误响应
	COMMON_FAIL = MsgCode{Code: -1, Msg: "错误"}

	//token过期或无效
	TOKEN_IS_EXPIRED = MsgCode{Code: -2, Msg: "Token过期或无效"}

	//参数错误
	PARAM_ERROR = MsgCode{Code: 1000, Msg: "参数错误"}

	//Query参数错误
	QUERY_PARAM_ERROR = MsgCode{Code: 1001, Msg: "Query参数错误"}

	//用户服务相关错误
	USER_ID_NOT_EXIST = MsgCode{Code: 2001, Msg: "用户ID不存在"}

	//websocket相关

	//用户未连接
	USER_NOT_CONNECTED = MsgCode{Code: 3001, Msg: "用户未连接"}
	//用户ID为空
	USER_ID_IS_EMPTY = MsgCode{Code: 3002, Msg: "用户ID为空"}
	//解析消息错误
	PARSE_MESSAGE_ERROR = MsgCode{Code: 3003, Msg: "解析消息错误"}
	//房间不存在
	ROOM_NOT_EXIST = MsgCode{Code: 3004, Msg: "房间不存在"}
)

func CustomError(err error) MsgCode {
	return MsgCode{Code: -3, Msg: err.Error()}
}
