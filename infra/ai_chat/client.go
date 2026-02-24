package ai_chat

import (
	"ai_interview/biz/ai_chat"
	"ai_interview/conf"
	"ai_interview/pkg/zlog"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
)

type ChatClient struct {
	Client *arkruntime.Client
}

var chatClient *ChatClient

func InitAiChar() {
	client := arkruntime.NewClientWithApiKey(conf.GetConfig().AiChat.ApiKey, arkruntime.WithBaseUrl(conf.GetConfig().AiChat.BaseUrl))
	if client == nil {
		panic("ai连接失败")
	}

	chatClient = &ChatClient{Client: client}

	zlog.Infof("ai连接成功")

}

func GetAiClient() ai_chat.IChatService {
	return chatClient
}
