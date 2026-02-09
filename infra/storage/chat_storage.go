package storage

import (
	"ai_interview/biz/entity"
	"ai_interview/biz/repo"
	"ai_interview/infra/database"
	"ai_interview/infra/storage/po"
	"ai_interview/pkg/error_msg"
	"ai_interview/pkg/zlog"
	"context"
	"encoding/json"
	"fmt"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

var cs *ChatStorage

type ChatStorage struct {
	db *gorm.DB
}

func InitChatStorage() {
	db := database.GetDB()
	if err := db.AutoMigrate(&po.Conversation{}, &po.Interview{}); err != nil {
		panic("chat部分表自动迁移失败: " + err.Error())
	}

	cs = &ChatStorage{db}
}

func GetChatStorage() repo.ChatRepo {
	return cs
}

// 创建并开始面试
func (c *ChatStorage) CreateInterview(ctx context.Context, userID, talentID, conversationID string) error {
	if talentID == "" {
		return error_msg.TALENT_ID_NOT_NULL
	} else if conversationID == "" {
		return error_msg.CONVERSATION_ID_NOT_NULL
	}

	conversationPo := po.Conversation{
		ConversationID: conversationID,
		Messages:       datatypes.JSON("[]"),
	}

	interviewPo := po.Interview{
		ConversationID:     conversationID,
		UserID:             userID,
		TalentID:           talentID,
		InterviewStartTime: time.Now(),
	}

	tx := c.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			zlog.Errorf("事务执行异常，已回滚")
		}
	}()

	err := tx.Model(&po.Conversation{}).WithContext(ctx).Create(&conversationPo).Error
	if err != nil {
		tx.Rollback()
		return errorDB(err)
	}

	err = tx.Model(&po.Interview{}).WithContext(ctx).Create(&interviewPo).Error
	if err != nil {
		tx.Rollback()
		return errorDB(err)
	}

	err = tx.Commit().Error
	if err != nil {
		return errorDB(err)
	}
	return nil
}

func (c *ChatStorage) EndInterview(ctx context.Context, conversationID string) error {
	if conversationID == "" {
		return error_msg.CONVERSATION_ID_NOT_NULL
	}

	updates := map[string]interface{}{
		"interview_end_time": time.Now(),
	}
	err := c.db.Model(&po.Interview{}).WithContext(ctx).Where("conversation_id = ?", conversationID).Updates(updates).Error
	if err != nil {
		return errorDB(err)
	}
	return nil
}

func (c *ChatStorage) GetInterviewsByUserID(ctx context.Context, userID string) ([]entity.Interview, error) {
	var interviews []entity.Interview
	err := c.db.Model(&po.Interview{}).WithContext(ctx).Where("user_id = ?", userID).Find(&interviews).Error
	if err != nil {
		return nil, errorDB(err)
	}
	return interviews, nil
}

func (c *ChatStorage) GetChatRecordByConversationID(ctx context.Context, conversationID string) ([]entity.Message, error) {
	if conversationID == "" {
		return nil, error_msg.CONVERSATION_ID_NOT_NULL
	}

	var messages []entity.Message
	var conversation po.Conversation
	err := c.db.Model(&po.Conversation{}).WithContext(ctx).Where("conversation_id = ?", conversationID).First(&conversation).Error
	if err != nil {
		return nil, errorDB(err)
	}
	err = json.Unmarshal(conversation.Messages, &messages)
	if err != nil {
		return nil, fmt.Errorf("反序列化消息失败: %v", err)
	}
	return messages, nil
}

func (c *ChatStorage) SaveAllMessages(ctx context.Context, conversationID string, messages []entity.Message) error {
	if conversationID == "" {
		return error_msg.CONVERSATION_ID_NOT_NULL
	}
	jsonData, err := json.Marshal(messages)
	if err != nil {
		return fmt.Errorf("序列化消息失败: %v", err)
	}
	updates := map[string]interface{}{
		"messages": datatypes.JSON(jsonData),
	}
	err = c.db.Model(&po.Conversation{}).WithContext(ctx).Where("conversation_id = ?", conversationID).Updates(updates).Error
	if err != nil {
		return errorDB(err)
	}
	return nil
}

func (c *ChatStorage) AddMessage(ctx context.Context, conversationID string, message entity.Message) error {
	if conversationID == "" {
		return error_msg.CONVERSATION_ID_NOT_NULL
	} else if message.Role == "" || (message.Role != "user" && message.Role != "assistant" && message.Role != "system") {
		return error_msg.MESSAGE_ROLE_ERROR
	} else if message.Content == "" {
		return error_msg.MESSAGE_CONTENT_NOT_NULL
	}
	//TODO 需要加锁解决并发问题 （同时更新可能会有并发问题）

	var conversation po.Conversation
	err := c.db.Model(&po.Conversation{}).WithContext(ctx).Where("conversation_id = ?", conversationID).First(&conversation).Error
	if err != nil {
		return errorDB(err)
	}
	var messages []entity.Message
	err = json.Unmarshal(conversation.Messages, &messages)
	if err != nil {
		return fmt.Errorf("反序列化消息失败: %v", err)
	}

	messages = append(messages, message)

	jsonData, err := json.Marshal(messages)
	if err != nil {
		return fmt.Errorf("序列化消息失败: %v", err)
	}
	updates := map[string]interface{}{
		"messages": datatypes.JSON(jsonData),
	}
	err = c.db.Model(&po.Conversation{}).WithContext(ctx).Where("conversation_id = ?", conversationID).Updates(updates).Error
	if err != nil {
		return errorDB(err)
	}
	return nil
}

func (c *ChatStorage) DeleteInterview(ctx context.Context, conversationID string) error {
	//TODO implement me
	panic("implement me")
}
