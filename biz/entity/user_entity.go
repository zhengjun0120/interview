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

type talentUserCtxKey struct{}

func WithUserID(ctx context.Context, userID string) context.Context {
	ctx = context.WithValue(ctx, userCtxKey{}, userID)
	return ctx
}

func WithTalentID(ctx context.Context, talentID string) context.Context {
	ctx = context.WithValue(ctx, talentUserCtxKey{}, talentID)
	return ctx
}

func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userCtxKey{}).(string)
	return userID, ok
}
func GetTalentID(ctx context.Context) (string, bool) {
	talentID, ok := ctx.Value(talentUserCtxKey{}).(string)
	return talentID, ok
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
