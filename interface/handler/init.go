package handler

import (
	"ai_interview/biz/types"
)

type Handler struct {
	UserServer       types.IUserService
	ResumeServer     types.IResumeService
	CodeServer       types.ICodeService
	JobProfileServer types.IJobProfileService
	ChatServer       types.IChatService
}

var handler *Handler

func GetHandler() *Handler {
	return handler
}

func InitHandler(userServer types.IUserService, resumeServer types.IResumeService, jobProfileServer types.IJobProfileService, codeServer types.ICodeService) {
	handler = &Handler{
		UserServer:       userServer,
		ResumeServer:     resumeServer,
		JobProfileServer: jobProfileServer,
		CodeServer:       codeServer,
	}
}
