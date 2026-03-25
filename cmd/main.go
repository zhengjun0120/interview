package main

import (
	"ai_interview/biz/code_service"
	"ai_interview/biz/job_profile_service"
	"ai_interview/biz/resume_service"
	"ai_interview/biz/user_service"
	"ai_interview/conf"
	"ai_interview/infra/ai_chat"
	"ai_interview/infra/cos"
	"ai_interview/infra/database"
	"ai_interview/infra/storage"
	"ai_interview/interface/handler"
	"ai_interview/interface/router"
	"ai_interview/pkg/zlog"
	"ai_interview/util"
)

func main() {
	// 初始化配置
	conf.InitConfig()

	//初始化日志
	zlog.InitLog("dev", "info")

	// 初始化数据库连接
	database.InitDB()

	//初始化redis连接
	database.InitRedis()

	//初始化雪花id
	util.InitSnowflake()

	//初始化ai
	ai_chat.InitAiChar()

	//初始化Cos
	cos.InitCos()

	//初始化仓库
	storage.InitUserStorage()
	storage.InitResumeStorage()
	storage.InitChatStorage()
	storage.InitJobProfileStorage()
	storage.InitCodeStorage()

	// 初始化服务
	us := user_service.NewUserService(storage.GetUserStorage(), storage.GetCodeStorage())
	cs := code_service.NewCodeService(storage.GetCodeStorage())
	rs := resume_service.NewResumeService(storage.GetResumeStorage(), storage.GetChatStorage(), storage.GetJobProfileStorage(), ai_chat.GetAiClient())
	js := job_profile_service.NewJobProfileService(storage.GetJobProfileStorage())

	// 初始化handler
	handler.InitHandler(us, rs, js, cs)

	// 启动服务
	router.RunServer()

}
