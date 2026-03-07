package types

import (
	"context"
	"github.com/gorilla/websocket"
)

type IChatService interface {
	// 创建面试房间
	CreateInterviewRoom(ctx context.Context, req *CreateInterviewRoomRequest) (*CreateInterviewRoomResponse, error)

	// 加入面试房间
	JoinInterviewRoom(ctx context.Context, req *JoinInterviewRoomRequest) error

	//CheckRoomPermission 预先检查是否符合连接ws条件
	CheckRoomPermission(ctx context.Context, req *CheckRoomPermissionRequest) (*CheckRoomPermissionResponse, error)

	//给hr的某个问题打标签
	AddTagToMessage(ctx context.Context, req *AddTagToMessageRequest) error

	// 获取面试时的实时追问
	GetAiSuggestion(ctx context.Context, req *GetAiSuggestionRequest) (*GetAiSuggestionResponse, error)

	//结束面试
	EndInterview(ctx context.Context, req *EndInterviewRequest) error
}

// 需要被面试者的ID
type CreateInterviewRoomRequest struct {
	TalentID string
}

// 返回房间ID 被面试者的临时token
type CreateInterviewRoomResponse struct {
	RoomID      string
	TalentToken string
}

type JoinInterviewRoomRequest struct {
	RoomID string
	Conn   *websocket.Conn
}

type CheckRoomPermissionRequest struct {
	RoomID string
}

type CheckRoomPermissionResponse struct {
	Success bool
}

type AddTagToMessageRequest struct {
	MessageID string
	RoomID    string
	Tag       string
}
type GetAiSuggestionRequest struct {
	RoomID string
}
type GetAiSuggestionResponse struct {
	Text string
}

type EndInterviewRequest struct {
	RoomID string
}
