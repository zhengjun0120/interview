package caster

import (
	"ai_interview/biz/types"
	"ai_interview/interface/def"
)

func CastLoginReq2ServiceParams(req *def.LoginReq) *types.LoginParams {
	return &types.LoginParams{
		Email:    req.Email,
		Password: req.Password,
	}
}

func CastServiceResp2LoginResp(resp *types.LoginResponse) *def.LoginResp {
	return &def.LoginResp{
		Token:    resp.Token,
		Username: resp.Username,
		Email:    resp.Email,
	}
}
