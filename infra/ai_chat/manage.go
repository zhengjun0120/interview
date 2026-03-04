package ai_chat

import (
	"ai_interview/biz/entity"
	"ai_interview/pkg/zlog"
	"encoding/json"
	"fmt"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
	"strings"
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

type InterviewMessageJson struct {
	Role string `json:"role"`
	Text string `json:"text"`
}

func buildInterviewMessage(message []entity.InterviewMessage, resumeStr string) (string, error) {
	var resp strings.Builder
	if resumeStr != "" {
		resp.WriteString("# 候选人完整简历 \n")
		resp.WriteString(resumeStr)
	}

	messages := make([]InterviewMessageJson, len(message))
	for i, msg := range message {
		messages[i] = InterviewMessageJson{
			Role: msg.From,
			Text: msg.Text,
		}
	}

	messagesStr, err := json.Marshal(messages)
	if err != nil {
		zlog.Errorf("构建面试消息失败 error: %v", err)
		return "", fmt.Errorf("构建面试消息失败 error: %w", err)
	}
	resp.WriteString("\n# 面试历史对话记录 \n")
	resp.WriteString(string(messagesStr))

	return resp.String(), nil
}
