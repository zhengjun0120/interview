package caster

import (
	"ai_interview/biz/types"
	"ai_interview/interface/def"
)

func CastCreateInterviewReq2ServiceParams(req *def.CreateInterviewRequest) *types.CreateInterviewRoomRequest {
	return &types.CreateInterviewRoomRequest{
		TalentID: req.TalentID,
	}
}

func CastServiceResp2CreateInterviewResp(resp *types.CreateInterviewRoomResponse) *def.CreateInterviewResponse {
	return &def.CreateInterviewResponse{
		RoomID:      resp.RoomID,
		TalentToken: resp.TalentToken,
	}
}

func CastCheckRoomPermissionReq2ServiceParams(req *def.CheckRoomPermissionRequest) *types.CheckRoomPermissionRequest {
	return &types.CheckRoomPermissionRequest{
		RoomID: req.RoomID,
	}
}

func CastServiceResp2CheckRoomPermissionResp(resp *types.CheckRoomPermissionResponse) *def.CheckRoomPermissionResponse {
	return &def.CheckRoomPermissionResponse{
		Success: resp.Success,
	}
}

func CastJoinInterviewRoomReq2ServiceParams(req *def.JoinInterviewRoomRequest) *types.JoinInterviewRoomRequest {
	return &types.JoinInterviewRoomRequest{
		RoomID: req.RoomID,
		Conn:   req.Conn,
	}
}

func CastAddTagToMessageReq2ServiceParams(req *def.AddTagToMessageRequest) *types.AddTagToMessageRequest {
	return &types.AddTagToMessageRequest{
		MessageID: req.MessageID,
		RoomID:    req.RoomID,
		Tag:       req.Tag,
	}
}

func CastGetAiSuggestionReq2ServiceParams(req *def.GetAiSuggestionRequest) *types.GetAiSuggestionRequest {
	return &types.GetAiSuggestionRequest{
		RoomID: req.RoomID,
	}
}
func CastServiceResp2GetAiSuggestionResp(resp *types.GetAiSuggestionResponse) *def.GetAiSuggestionResponse {
	return &def.GetAiSuggestionResponse{
		Text: resp.Text,
	}
}

func CastEndInterviewReq2ServiceParams(req *def.EndInterviewRequest) *types.EndInterviewRequest {
	return &types.EndInterviewRequest{
		RoomID: req.RoomID,
	}
}
