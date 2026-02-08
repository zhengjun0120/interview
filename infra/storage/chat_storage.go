package storage

import (
	"ai_interview/infra/database"
	"ai_interview/infra/storage/po"
	"gorm.io/gorm"
)

var ch *ChatStorage

type ChatStorage struct {
	db *gorm.DB
}

func InitChatStorage() {
	db := database.GetDB()
	if err := db.AutoMigrate(&po.Conversation{}, &po.Interview{}); err != nil {
		panic("chat部分表自动迁移失败: " + err.Error())
	}

	ch = &ChatStorage{db}
}

func GetChatStorage() *ChatStorage {
	return ch
}
