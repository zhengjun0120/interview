package ai_chat

import (
	"ai_interview/biz/entity"
	"ai_interview/conf"
	"context"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
)

func Chat(ctx context.Context, messages []entity.Message) (string, error) {

	reqMessage := Message2ReqMessage(messages)

	req := model.CreateChatCompletionRequest{
		Model:    conf.GetConfig().AiChat.ModelId,
		Messages: reqMessage,
		Thinking: &model.Thinking{
			Type: model.ThinkingTypeAuto,
		},
	}

	resp, err := client.CreateChatCompletion(ctx, req)

	if err != nil {
		return "", err
	}
	return *resp.Choices[0].Message.Content.StringValue, nil

}
