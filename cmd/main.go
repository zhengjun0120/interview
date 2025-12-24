package main

import (
	"ai_interview/biz/user_service"
	"ai_interview/conf"
	"ai_interview/infra/database"
	"ai_interview/interface/handler"
	"ai_interview/interface/router"
)

func main() {
	conf.InitConfig()

	// 初始化数据库连接
	database.InitDB()

	// 初始化服务
	us := user_service.NewUserService()

	// 初始化handler
	handler.InitHandler(us)

	// 启动服务
	router.RunServer()

}
