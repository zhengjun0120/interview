package entity

import "context"

type User struct {
	UserID   string
	Username string
	Password string
	Email    string
	Type     string
}

type userCtxKey struct{}

func WithUserID(ctx context.Context, userID string) context.Context {
	ctx = context.WithValue(ctx, userCtxKey{}, userID)
	return ctx
}

func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userCtxKey{}).(string)
	return userID, ok
}

func NewUser(userID, username, password, email, userType string) *User {

	return &User{
		UserID:   userID,
		Username: username,
		Password: password,
		Email:    email,
		Type:     userType,
	}
}
