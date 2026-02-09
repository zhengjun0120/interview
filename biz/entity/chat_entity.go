package entity

import "time"

type Message struct {
	Role      string    `json:"role"`       // 角色
	Content   string    `json:"content"`    // 文本
	CreatedAt time.Time `json:"created_at"` // 创建时间
}

// Conversation ai 聊天会话
type Conversation struct {
	ConversationID string    // 聊天会话ID
	Messages       []Message // 聊天消息
}

// 面试
type Interview struct {
	UserID             string    // 用户ID
	InterviewID        string    // 面试ID
	TalentID           string    // 人才ID
	ConversationID     string    // 聊天会话ID
	InterviewStartTime time.Time // 面试开始时间
	InterviewEndTime   time.Time // 面试结束时间
	CreatedAt          time.Time // 创建时间
}
