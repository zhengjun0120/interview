package def

import "github.com/gorilla/websocket"

type CreateInterviewRequest struct {
	TalentID string `json:"talent_id"`
}

type CreateInterviewResponse struct {
	RoomID      string `json:"room_id"`
	TalentToken string `json:"talent_token"`
}

type CheckRoomPermissionRequest struct {
	RoomID string `json:"room_id"`
}

type CheckRoomPermissionResponse struct {
	Success bool `json:"success"`
}

type JoinInterviewRoomRequest struct {
	RoomID string
	Conn   *websocket.Conn
}
