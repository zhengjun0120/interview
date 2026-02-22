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
)
