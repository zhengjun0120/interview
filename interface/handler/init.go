package handler

import "ai_interview/biz/types"

type Handler struct {
	UserServer types.IUserService
}

var handler Handler

func GetHandler() Handler {
	return handler
}

func InitHandler(userServer types.IUserService) {
	handler.UserServer = userServer
}
