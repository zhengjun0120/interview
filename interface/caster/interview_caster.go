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
