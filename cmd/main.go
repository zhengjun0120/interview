package main

import (
	"ai_interview/conf"
	"ai_interview/infra/database"
)

func main() {
	conf.InitConfig()

	// 初始化数据库连接
	database.InitDB()

}
