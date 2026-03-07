package storage

import (
	"ai_interview/biz/chat_service/ws"
	"ai_interview/biz/entity"
	"ai_interview/biz/repo"
	"ai_interview/infra/database"
	"ai_interview/infra/storage/po"
	"ai_interview/pkg/error_msg"
	"context"
	"errors"
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
	if err := db.AutoMigrate(&po.Interview{}, &po.InterviewMessage{}); err != nil {
		panic("chat部分表自动迁移失败: " + err.Error())
	}

	cs = &ChatStorage{db}
}

func GetChatStorage() repo.ChatRepo {
	return cs
}

func (c *ChatStorage) po2EntityInterview(po po.Interview) entity.Interview {
	return entity.Interview{
		RoomID:             po.RoomID,
		UserID:             po.UserID,
		TalentID:           po.TalentID,
		TalentReport:       po.TalentReport,
		InterviewStartTime: po.InterviewStartTime,
		InterviewEndTime:   po.InterviewEndTime,
		CreatedAt:          po.CreatedAt,
	}
}

// 创建并开始面试
func (c *ChatStorage) CreateInterview(ctx context.Context, userID, talentID, roomID string) error {
	if talentID == "" {
		return error_msg.TALENT_ID_NOT_NULL
	} else if roomID == "" {
		return error_msg.ROOM_ID_NOT_NULL
	}

	interviewPo := po.Interview{
		RoomID:             roomID,
		UserID:             userID,
		TalentID:           talentID,
		InterviewStartTime: time.Now(),
	}

	err := c.db.Model(&po.Interview{}).WithContext(ctx).Create(&interviewPo).Error
	if err != nil {
		return errorDB(err)
	}

	return nil
}

func (c *ChatStorage) EndInterview(ctx context.Context, roomID, userID string) error {
	if roomID == "" {
		return error_msg.ROOM_ID_NOT_NULL
	}

	// 查询面试记录
	var interview po.Interview
	err := c.db.Model(&po.Interview{}).WithContext(ctx).Where("room_id = ? AND user_id = ?", roomID, userID).First(&interview).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return error_msg.ROOM_NOT_EXIST
		}
		return errorDB(err)
	}

	// 更新面试记录
	// 开启事务 更新面试时间和面试状态
	err = c.db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		updates := map[string]interface{}{
			"interview_end_time": &now,
		}
		err = tx.Model(&po.Interview{}).WithContext(ctx).Where("room_id = ? AND user_id = ?", roomID, userID).Updates(updates).Error
		if err != nil {
			return errorDB(err)
		}

		err = tx.Model(&po.TalentPool{}).WithContext(ctx).Where("talent_id = ?", interview.TalentID).Update("interview_status", entity.InterviewStatusInterviewed).Error
		if err != nil {
			return errorDB(err)
		}
		return nil
	})

	if err != nil {
		return errorDB(err)
	}

	return nil
}

func (c *ChatStorage) GetInterviewsByUserID(ctx context.Context, userID string) ([]entity.Interview, error) {
	var interviews []po.Interview
	err := c.db.Model(&po.Interview{}).WithContext(ctx).Where("user_id = ?", userID).Find(&interviews).Error
	if err != nil {
		return nil, errorDB(err)
	}

	resp := make([]entity.Interview, len(interviews))
	for i, v := range interviews {

		resp[i] = c.po2EntityInterview(v)
	}

	return resp, nil
}

func (c *ChatStorage) SaveInterviewMessage(ctx context.Context, message entity.InterviewMessage) error {
	if message.MessageID == "" {
		return error_msg.MESSAGE_ID_NOT_NULL
	} else if message.From != ws.HR && message.From != ws.CANDIDATE {
		return error_msg.MESSAGE_FROM_ERROR
	} else if message.RoomID == "" {
		return error_msg.ROOM_ID_NOT_NULL
	} else if message.Text == "" {
		return error_msg.MESSAGE_TEXT_NOT_NULL
	}

	var messagePo = po.InterviewMessage{
		MessageID: message.MessageID,
		UserID:    message.UserID,
		From:      message.From,
		Tag:       message.Tag,
		Text:      message.Text,
		RoomID:    message.RoomID,
		Type:      message.Type,
	}
	err := c.db.Model(&po.InterviewMessage{}).WithContext(ctx).Create(&messagePo).Error
	if err != nil {
		return errorDB(err)
	}
	return nil
}

func (c *ChatStorage) GetInterviewMessageByRoomID(ctx context.Context, roomID string) ([]entity.InterviewMessage, error) {
	if roomID == "" {
		return nil, error_msg.ROOM_ID_NOT_NULL
	}

	var messages []po.InterviewMessage
	err := c.db.Model(&po.InterviewMessage{}).WithContext(ctx).Where("room_id= ?", roomID).Order("created_at ASC").Find(&messages).Error
	if err != nil {
		return nil, errorDB(err)
	}

	resp := make([]entity.InterviewMessage, len(messages))
	for i, v := range messages {
		resp[i] = entity.InterviewMessage{
			MessageID: v.MessageID,
			UserID:    v.UserID,
			From:      v.From,
			Tag:       v.Tag,
			Text:      v.Text,
			RoomID:    v.RoomID,
			Type:      v.Type,
		}
	}

	return resp, nil
}

func (c *ChatStorage) AddTagToMessage(ctx context.Context, messageID string, tag string, roomID string) error {
	if messageID == "" {
		return error_msg.MESSAGE_ID_NOT_NULL
	} else if tag == "" {
		return error_msg.TAG_NOT_NULL
	}

	var check po.InterviewMessage
	err := c.db.Model(&po.InterviewMessage{}).WithContext(ctx).Where("message_id = ? AND room_id = ?", messageID, roomID).First(&check).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return error_msg.MESSAGE_NOT_EXIST
		}
		return errorDB(err)
	}

	err = c.db.Model(&po.InterviewMessage{}).WithContext(ctx).Where("message_id = ?", messageID).Update("tag", tag).Error
	if err != nil {
		return errorDB(err)
	}
	return nil
}

func (c *ChatStorage) CheckRoom(ctx context.Context, roomID, userID string) (bool, error) {
	var interview po.Interview
	err := c.db.Model(&po.Interview{}).WithContext(ctx).Where("room_id = ? ", roomID).First(&interview).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, errorDB(err)
	}
	if interview.UserID != userID && interview.TalentID != userID {
		return false, error_msg.ROOM_NOT_EXIST
	}
	if interview.InterviewEndTime != nil {
		return false, error_msg.ROOM_ALREADY_ENDED
	}
	return true, nil
}

func (c *ChatStorage) GetInterviewByRoomID(ctx context.Context, roomID, userID string) (*entity.Interview, error) {
	if roomID == "" {
		return nil, error_msg.ROOM_ID_NOT_NULL
	}
	var interview po.Interview
	err := c.db.Model(&po.Interview{}).WithContext(ctx).Where("room_id = ? AND user_id = ?", roomID, userID).First(&interview).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, error_msg.ROOM_NOT_EXIST
		}
		return nil, errorDB(err)
	}
	resp := c.po2EntityInterview(interview)
	return &resp, nil
}

func (c *ChatStorage) GetInterviewByTalentID(ctx context.Context, talentID, userID string) (*entity.Interview, error) {
	if talentID == "" {
		return nil, error_msg.TALENT_ID_NOT_NULL
	}

	var interview po.Interview
	err := c.db.Model(&po.Interview{}).WithContext(ctx).Where("talent_id = ? AND user_id = ?", talentID, userID).First(&interview).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, error_msg.ROOM_NOT_EXIST
		}
		return nil, errorDB(err)
	}

	resp := c.po2EntityInterview(interview)

	return &resp, nil

}

func (c *ChatStorage) SaveTalentReport(ctx context.Context, report datatypes.JSON, userID, talentID string) error {
	if talentID == "" {
		return error_msg.TALENT_ID_NOT_NULL
	}
	updates := map[string]interface{}{
		"talent_report": report,
	}

	err := c.db.Model(&po.Interview{}).WithContext(ctx).Where("talent_id = ? AND user_id = ?", talentID, userID).Updates(updates).Error
	if err != nil {
		return errorDB(err)
	}
	return nil
}
