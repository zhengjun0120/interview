package user_service

import (
	"ai_interview/biz/entity"
	"ai_interview/biz/repo"
	"ai_interview/biz/types"
	"ai_interview/pkg/error_msg"
	"ai_interview/util"
	"context"
)

var (
	HR     = "hr"
	Seeker = "seeker"
)

type UserService struct {
	// TODO: 注入需要的依赖
	userRepo repo.UserRepo
	codeRepo repo.CodeRepo
}

func NewUserService(userRepo repo.UserRepo, codeRepo repo.CodeRepo) *UserService {
	return &UserService{
		userRepo: userRepo,
		codeRepo: codeRepo,
	}
}

func (u *UserService) Login(ctx context.Context, req *types.LoginParams) (*types.LoginResponse, error) {

	user, err := u.userRepo.GetUserByEmail(ctx, req.Email, req.Type)
	if err != nil {
		return nil, err
	}

	if !util.CheckPassword(user.Password, req.Password) {
		return nil, error_msg.PASSWORD_ERROR
	}

	token, err := util.GenerateJWT(user.UserID)
	if err != nil {
		return nil, err
	}

	return &types.LoginResponse{
		Token:    token,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

// 注册
func (u *UserService) Register(ctx context.Context, req *types.RegisterParams) (*types.RegisterResponse, error) {
	if req.Type != HR && req.Type != Seeker {
		return nil, error_msg.USER_TYPE_ERROR
	} else if req.Password == "" {
		return nil, error_msg.PASSWORD_NOT_NULL
	} else if req.Username == "" {
		return nil, error_msg.USERNAME_NOT_NULL
	}

	if err := util.CheckEmail(req.Email); err != nil {
		return nil, error_msg.Email_FORMAT_INVALID
	}

	ok, err := u.userRepo.CheckEmail(ctx, req.Email, req.Type)
	if err != nil {
		return nil, err
	}
	if ok {
		return nil, error_msg.EMAIL_EXIST
	}

	// 检查验证码
	if err := u.codeRepo.CaptchaCheck(ctx, types.CaptchaWayTypeRegister, req.Email, req.Code); err != nil {
		return nil, err
	}

	userID := util.GenerateStringID()
	hashPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		UserID:   userID,
		Username: req.Username,
		Password: hashPassword,
		Email:    req.Email,
		Type:     req.Type,
	}

	err = u.userRepo.CreateUser(ctx, user, req.Type)
	if err != nil {
		return nil, err
	}

	token, err := util.GenerateJWT(userID)
	if err != nil {
		return nil, err
	}

	return &types.RegisterResponse{
		Token:    token,
		Username: user.Username,
		Email:    user.Email,
	}, nil

}
