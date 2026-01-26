package user_service

import (
	"ai_interview/biz/entity"
	"ai_interview/biz/repo"
	"ai_interview/biz/types"
	"ai_interview/util"
	"context"
	"errors"
)

var (
	HR     = "hr"
	Seeker = "seeker"

	USER_ID_NOT_EXIST = errors.New("用户ID不存在")
	USER_NOT_EXIST    = errors.New("用户不存在")
	USER_TYPE_ERROR   = errors.New("用户类型错误")
	PASSWORD_NOT_NULL = errors.New("密码不能为空")
	USERNAME_NOT_NULL = errors.New("用户名不能为空")
	PASSWORD_ERROR    = errors.New("密码错误或用户不存在")

	EMAIL_EXIST     = errors.New("邮箱已注册")
	EMAIL_NOT_EXIST = errors.New("邮箱未注册")
)

type UserService struct {
	// TODO: 注入需要的依赖
	userRepo repo.UserRepo
}

func NewUserService(userRepo repo.UserRepo) *UserService {
	return &UserService{userRepo: userRepo}
}

func (u *UserService) Login(ctx context.Context, req *types.LoginParams) (*types.LoginResponse, error) {

	user, err := u.userRepo.GetUserByEmail(ctx, req.Email, req.Type)
	if err != nil {
		return nil, err
	}

	if !util.CheckPassword(user.Password, req.Password) {
		return nil, PASSWORD_ERROR
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
		return nil, USER_TYPE_ERROR
	} else if req.Password == "" {
		return nil, PASSWORD_NOT_NULL
	} else if req.Username == "" {
		return nil, USERNAME_NOT_NULL
	}

	ok, err := u.userRepo.CheckEmail(ctx, req.Email, req.Type)
	if err != nil {
		return nil, err
	}
	if ok {
		return nil, EMAIL_EXIST
	}

	// TODO: 检查验证码

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
