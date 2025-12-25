package types

import "context"

type IUserService interface {
	Login(ctx context.Context, req *LoginParams) (*LoginResponse, error)
	Register(ctx context.Context, req *RegisterParams) (*RegisterResponse, error)
}

type LoginParams struct {
	Email    string
	Password string
	Type     string
}

type LoginResponse struct {
	Token    string
	Username string
	Email    string
}

type RegisterParams struct {
	Email    string
	Password string
	Username string
	Type     string
	Code     string
}

type RegisterResponse struct {
	Token    string
	Username string
	Email    string
}
