package user_service

import (
	"ai_interview/biz/types"
	"context"
	"errors"
)

var (
	USER_ID_NOT_EXIST = errors.New("用户ID不存在")
)

type UserService struct {
	// TODO: 注入需要的依赖
}

func NewUserService() *UserService {
	return &UserService{}
}

func (u *UserService) Login(ctx context.Context, req *types.LoginParams) (*types.LoginResponse, error) {
	// TODO: 实现登录逻辑
	return nil, nil
}
