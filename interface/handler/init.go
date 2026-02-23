package handler

import (
	"ai_interview/biz/types"
)

type Handler struct {
	UserServer   types.IUserService
	ResumeServer types.IResumeService
	CodeServer   types.ICodeService
}

var handler Handler

func GetHandler() Handler {
	return handler
}

func InitHandler(userServer types.IUserService, resumeServer types.IResumeService) {
	handler.UserServer = userServer
	handler.ResumeServer = resumeServer
}
