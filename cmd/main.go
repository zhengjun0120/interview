package main

import (
	"ai_interview/biz/service/resume_service"
	"ai_interview/biz/service/user_service"
	"ai_interview/conf"
	"ai_interview/infra/ai_chat"
	"ai_interview/infra/cos"
	"ai_interview/infra/database"
	"ai_interview/infra/storage"
	"ai_interview/interface/handler"
	"ai_interview/interface/router"
	"ai_interview/util"
)

func main() {
	// 初始化配置
	conf.InitConfig()

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

	// 初始化服务
	us := user_service.NewUserService(storage.GetUserStorage())
	rs := resume_service.NewResumeService(storage.GetResumeStorage())

	// 初始化handler
	handler.InitHandler(us, rs)

	// 启动服务
	router.RunServer()

}
