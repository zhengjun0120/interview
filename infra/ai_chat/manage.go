package ai_chat

import (
	"ai_interview/biz/entity"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
)

func Message2ReqMessage(messages []entity.Message) []*model.ChatCompletionMessage {
	messageReq := make([]*model.ChatCompletionMessage, len(messages))

	for i, message := range messages {
		messageReq[i] = &model.ChatCompletionMessage{
			Role: message.Role,
			Content: &model.ChatCompletionMessageContent{
				StringValue: volcengine.String(message.Content),
			},
		}
	}

	return messageReq
}
