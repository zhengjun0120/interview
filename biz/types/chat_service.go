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
