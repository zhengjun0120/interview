package caster

import (
	"ai_interview/biz/types"
	"ai_interview/interface/def"
)

func CastLoginReq2ServiceParams(req *def.LoginReq) *types.LoginParams {
	return &types.LoginParams{
		Email:    req.Email,
		Password: req.Password,
		Type:     req.Type,
	}
}

func CastServiceResp2LoginResp(resp *types.LoginResponse) *def.LoginResp {
	return &def.LoginResp{
		Token:    resp.Token,
		Username: resp.Username,
		Email:    resp.Email,
	}
}

func CastRegisterReq2ServiceParams(req *def.RegisterReq) *types.RegisterParams {
	return &types.RegisterParams{
		Email:    req.Email,
		Password: req.Password,
		Username: req.Username,
		Type:     req.Type,
		Code:     req.Code,
	}
}

func CastServiceResp2RegisterResp(resp *types.RegisterResponse) *def.RegisterResp {
	return &def.RegisterResp{
		Token:    resp.Token,
		Username: resp.Username,
		Email:    resp.Email,
	}
}
