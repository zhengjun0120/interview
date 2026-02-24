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

	resumeGroup := api.Group("/resume", middleware.JWTAuth())
	loadResumeService(resumeGroup)

	talentGroup := api.Group("/talent", middleware.JWTAuth())
	loadTalentService(talentGroup)

	jobProfileGroup := api.Group("/job_profile", middleware.JWTAuth())
	loadJobProfileService(jobProfileGroup)

	return r

}

func loadUserService(r *gin.RouterGroup) {
	r.POST("/login", handler.GetHandler().Login)
}
func loadResumeService(r *gin.RouterGroup) {
	// /api/v1/resume/upload POST 上传简历
	r.POST("/upload", handler.GetHandler().UploadResume)

}

func loadTalentService(r *gin.RouterGroup) {
	// /api/v1/talent/get GET 获取 talent 列表
	r.GET("/get", handler.GetHandler().GetTalent)
}
func loadJobProfileService(r *gin.RouterGroup) {
	// /api/v1/job_profile/get GET 获取人才画像
	r.GET("/get", handler.GetHandler().GetJobProfile)

	// /api/v1/job_profile/save POST 保存人才画像
	r.POST("/save", handler.GetHandler().SaveJobProfile)
}
