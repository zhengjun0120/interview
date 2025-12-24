package entity

import "context"

type user struct {
	UserID   string
	UserName string
	Password string

	Email string

	Type string
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
