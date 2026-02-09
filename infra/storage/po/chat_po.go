package po

import (
	"gorm.io/datatypes"
	"time"
)

// Conversation 会话表
type Conversation struct {
	ID             int            `gorm:"column:id;primary_key;autoIncrement"`    //主键
	ConversationID string         `gorm:"column:conversation_id;not null;unique"` //会话ID
	Messages       datatypes.JSON `gorm:"column:messages"`                        //消息
	CreatedAt      time.Time      `gorm:"column:created_at"`                      //创建时间
	UpdatedAt      time.Time      `gorm:"column:updated_at"`                      //更新时间
}

type Interview struct {
	ID                 int       `gorm:"column:id;primary_key;autoIncrement"`    //主键
	ConversationID     string    `gorm:"column:conversation_id;not null;unique"` //会话ID
	UserID             string    `gorm:"column:user_id;not null"`                //用户ID
	TalentID           string    `gorm:"column:talent_id;not null"`              //人才ID
	InterviewStartTime time.Time `gorm:"column:interview_start_time;not null"`   //面试开始时间
	InterviewEndTime   time.Time `gorm:"column:interview_end_time"`              //面试结束时间

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}
