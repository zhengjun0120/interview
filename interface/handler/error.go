package handler

import (
	"ai_interview/biz/user_service"
	"ai_interview/pkg/response"
	"errors"
)

func ErrorToMsgCode(err error) response.MsgCode {
	if err == nil {
		return response.SUCCESS
	}

	if errors.Is(err, user_service.USER_ID_NOT_EXIST) {
		return response.USER_ID_NOT_EXIST
	}

	return response.COMMON_FAIL
}
