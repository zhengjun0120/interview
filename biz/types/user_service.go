package types

import "context"

type IUserService interface {
	Login(ctx context.Context, req *LoginParams) (*LoginResponse, error)
}

type LoginParams struct {
	Email    string
	Password string
}

type LoginResponse struct {
	Token    string
	Username string
	Email    string
}
