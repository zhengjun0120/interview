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

// Interview 面试表
type Interview struct {
	ID                 int        `gorm:"column:id;primary_key;autoIncrement"`  //主键
	RoomID             string     `gorm:"column:room_id;not null;unique"`       //房间ID
	UserID             string     `gorm:"column:user_id;not null"`              //用户ID
	TalentID           string     `gorm:"column:talent_id;not null"`            //人才ID
	InterviewStartTime time.Time  `gorm:"column:interview_start_time;not null"` //面试开始时间
	InterviewEndTime   *time.Time `gorm:"column:interview_end_time"`            //面试结束时间

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

type InterviewMessage struct {
	ID        int    `gorm:"column:id;primary_key;autoIncrement"`
	MessageID string `gorm:"column:message_id;not null;unique"`
	RoomID    string `gorm:"column:room_id;not null"`
	From      string `gorm:"column:from;not null"`       // 发送者Hr 或者 candidate
	Text      string `gorm:"column:text;not null"`       // 消息内容
	Tag       string `gorm:"column:tag"`                 // HR标签
	Type      string `gorm:"column:type;default:'chat'"` // 消息类型

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}
