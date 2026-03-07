package def

import "github.com/gorilla/websocket"

// 创建面试房间请求
type CreateInterviewRequest struct {
	TalentID string `json:"talent_id"`
}

type CreateInterviewResponse struct {
	RoomID      string `json:"room_id"`
	TalentToken string `json:"talent_token"`
}

// 检查房间权限请求
type CheckRoomPermissionRequest struct {
	RoomID string `json:"room_id"`
}

type CheckRoomPermissionResponse struct {
	Success bool `json:"success"`
}

// 加入面试房间请求
type JoinInterviewRoomRequest struct {
	RoomID string
	Conn   *websocket.Conn
}

// 添加标签请求
type AddTagToMessageRequest struct {
	MessageID string `json:"message_id"`
	RoomID    string `json:"room_id"`
	Tag       string `json:"tag"`
}

// 获取AI建议请求
type GetAiSuggestionRequest struct {
	RoomID string `json:"room_id"`
}

type GetAiSuggestionResponse struct {
	Text string `json:"text"`
}

// 结束面试请求
type EndInterviewRequest struct {
	RoomID string `json:"room_id"`
}
