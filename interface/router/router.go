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

	interviewGroup := api.Group("/interview", middleware.JWTAuth())
	loadInterviewService(interviewGroup)

	return r

}

func loadUserService(r *gin.RouterGroup) {
	r.POST("/login", handler.GetHandler().Login)
}
func loadResumeService(r *gin.RouterGroup) {
	// /api/v1/resume/upload POST 上传简历
	r.POST("/upload", handler.GetHandler().UploadResume)

	// /api/v1/resume/get_url?talent_id=123 GET 获取简历 url
	r.GET("/get_url", handler.GetHandler().GetResumeUrl)

}
func loadInterviewService(r *gin.RouterGroup) {
	// /api/v1/interview/create POST 创建面试
	r.POST("/create", handler.GetHandler().CreateInterview)

	// /api/v1/interview/check_room_permission POST 检查房间权限
	r.POST("/check_room_permission", handler.GetHandler().CheckRoomPermission)

	// /api/v1/interview/ws/join?room_id=123&token=123 GET 加入面试房间
	r.GET("/ws/join", handler.GetHandler().JoinInterviewRoom)

	// /api/v1/interview/message/add_tag POST 添加消息标签
	r.POST("/message/add_tag", handler.GetHandler().AddTagToMessage)

	// /api/v1/interview/message/ai_suggestion POST 获取消息的 AI 建议
	r.POST("/message/ai_suggestion", handler.GetHandler().GetAiSuggestion)
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
