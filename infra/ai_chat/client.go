package ai_chat

import (
	"ai_interview/conf"
	"ai_interview/pkg/zlog"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
)

var client *arkruntime.Client

func InitAiChar() {
	client = arkruntime.NewClientWithApiKey(conf.GetConfig().AiChat.ApiKey, arkruntime.WithBaseUrl(conf.GetConfig().AiChat.BaseUrl))
	if client == nil {
		panic("ai连接失败")
	}

	zlog.Infof("ai连接成功")

}

func GetAiClient() *arkruntime.Client {
	return client
}
