package router

import (
	"ai_interview/conf"
	"ai_interview/interface/handler"
	"ai_interview/interface/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
)

func RunServer() {
	r := register()

	// 启动服务
	port := conf.GetConfig().Server.Port
	host := conf.GetConfig().Server.Host
	r.Run(host + ":" + port)
	fmt.Println("服务启动成功")
}

func register() *gin.Engine {
	r := gin.Default()
	//跨域中间件
	r.Use(middleware.CorsMiddleware())

	api := r.Group("/api/v1")

	// 用户相关接口 不需要鉴权
	userGroup := api.Group("/user")
	loadUserService(userGroup)

	return r

}

func loadUserService(r *gin.RouterGroup) {
	r.POST("/login", handler.GetHandler().Login)
}
