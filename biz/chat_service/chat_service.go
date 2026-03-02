package chat_service

import (
	"ai_interview/biz/chat_service/ws"
	"ai_interview/biz/entity"
	"ai_interview/biz/repo"
	"ai_interview/biz/types"
	"ai_interview/pkg/error_msg"
	"ai_interview/util"
	"context"
)

type ChatService struct {
	ChatRepo   repo.ChatRepo
	ResumeRepo repo.ResumeRepo
}

func NewChatService(chatRepo repo.ChatRepo, resumeRepo repo.ResumeRepo) *ChatService {
	return &ChatService{ChatRepo: chatRepo, ResumeRepo: resumeRepo}
}

func (c *ChatService) CreateInterviewRoom(ctx context.Context, req *types.CreateInterviewRoomRequest) (*types.CreateInterviewRoomResponse, error) {
	userID, ok := entity.GetUserID(ctx)
	if !ok {
		return nil, error_msg.GET_USER_ID_ERROR
	}

	ok, err := c.ResumeRepo.CheckTalentID(ctx, req.TalentID)
	if !ok {
		if err != nil {
			return nil, err
		} else {
			return nil, error_msg.TALENT_ID_NOT_EXIST
		}
	}

	roomID := util.GenerateStringID()
	err = c.ChatRepo.CreateInterview(ctx, userID, req.TalentID, roomID)
	if err != nil {
		return nil, err
	}

	talentToken, err := util.GenerateJWT(req.TalentID)
	if err != nil {
		return nil, err
	}
	return &types.CreateInterviewRoomResponse{RoomID: roomID, TalentToken: talentToken}, nil
}

func (c *ChatService) CheckRoomPermission(ctx context.Context, req *types.CheckRoomPermissionRequest) (*types.CheckRoomPermissionResponse, error) {
	var resp = &types.CheckRoomPermissionResponse{Success: false}
	var userID string
	userID, ok := entity.GetUserID(ctx)
	if !ok {
		userID, ok = entity.GetTalentID(ctx)
		if !ok {
			return resp, error_msg.GET_USER_ID_ERROR
		}
	}

	//检查房间权限
	ok, err := c.ChatRepo.CheckRoom(ctx, req.RoomID, userID)
	if !ok {
		if err != nil {
			return resp, err
		}
		return resp, error_msg.ROOM_NOT_EXIST
	}

	resp.Success = true

	return resp, nil
}

func (c *ChatService) JoinInterviewRoom(ctx context.Context, req *types.JoinInterviewRoomRequest) error {
	var client = &ws.Client{
		RoomID: req.RoomID,
	}

	var userID string
	userID, ok := entity.GetUserID(ctx)
	client.Role = ws.HR
	if !ok {
		userID, ok = entity.GetTalentID(ctx)
		if !ok {
			return error_msg.GET_USER_ID_ERROR
		}
		client.Role = ws.CANDIDATE
	}
	client.UserID = userID

	//检查房间是否存在且有权限
	ok, err := c.ChatRepo.CheckRoom(ctx, req.RoomID, userID)
	if !ok {
		if err != nil {
			return err
		}
		return error_msg.ROOM_NOT_EXIST
	}

	client.Conn = req.Conn
	client.Send = make(chan []byte, 256)

	ws.GetWebSocketHub().Register <- client

	return nil

}
