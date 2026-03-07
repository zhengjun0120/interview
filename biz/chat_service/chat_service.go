package chat_service

import (
	"ai_interview/biz/ai_chat"
	"ai_interview/biz/chat_service/ws"
	"ai_interview/biz/entity"
	"ai_interview/biz/repo"
	"ai_interview/biz/types"
	"ai_interview/pkg/error_msg"
	"ai_interview/util"
	"context"
)

type ChatService struct {
	ChatRepo      repo.ChatRepo
	ResumeRepo    repo.ResumeRepo
	AiChatService ai_chat.IChatService
}

func NewChatService(chatRepo repo.ChatRepo, resumeRepo repo.ResumeRepo, aiChatService ai_chat.IChatService) *ChatService {
	return &ChatService{ChatRepo: chatRepo, ResumeRepo: resumeRepo, AiChatService: aiChatService}
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

func (c *ChatService) AddTagToMessage(ctx context.Context, req *types.AddTagToMessageRequest) error {
	userID, ok := entity.GetUserID(ctx)
	if !ok {
		return error_msg.GET_USER_ID_ERROR
	}

	_, err := c.ChatRepo.GetInterviewByRoomID(ctx, req.RoomID, userID)
	if err != nil {
		return err
	}
	err = c.ChatRepo.AddTagToMessage(ctx, req.MessageID, req.Tag, req.RoomID)
	if err != nil {
		return err
	}
	return nil
}

func (c *ChatService) GetAiSuggestion(ctx context.Context, req *types.GetAiSuggestionRequest) (*types.GetAiSuggestionResponse, error) {
	userID, ok := entity.GetUserID(ctx)
	if !ok {
		return nil, error_msg.GET_USER_ID_ERROR
	}

	interview, err := c.ChatRepo.GetInterviewByRoomID(ctx, req.RoomID, userID)
	if err != nil {
		return nil, err
	}

	resume, err := c.ResumeRepo.GetResumeByTalentID(ctx, interview.TalentID)
	if err != nil {
		return nil, err
	}

	messages, err := c.ChatRepo.GetInterviewMessageByRoomID(ctx, req.RoomID)
	if err != nil {
		return nil, err
	}

	aiSuggestion, err := c.AiChatService.InterviewAiSuggestionChat(ctx, messages, resume.ResumeStr, resume.ResumeUrl)
	if err != nil {
		return nil, err
	}
	return &types.GetAiSuggestionResponse{Text: aiSuggestion}, nil
}
