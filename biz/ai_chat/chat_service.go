package ai_chat

import (
	"ai_interview/biz/entity"
	"context"
)

type IChatService interface {
	Chat(ctx context.Context, messages []entity.Message) (string, error)
	DocxChat(ctx context.Context, dockUrl, jobProfile string) (string, error)
}
