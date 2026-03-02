package ai_chat

import (
	"ai_interview/biz/entity"
	"context"
)

type IChatService interface {
	Chat(ctx context.Context, messages []entity.Message) (string, error)
	DocxToTalentDataChat(ctx context.Context, dockUrl, jobProfile string) (string, error)
	DocxToResumeStrChat(ctx context.Context, dockUrl string) (string, error)
}
