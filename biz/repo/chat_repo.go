package repo

import (
	"ai_interview/biz/entity"
	"context"
)

type ChatRepo interface {
	//新建面试
	CreateInterview(ctx context.Context, userID, talentID, conversationID string) error

	//结束面试
	EndInterview(ctx context.Context, conversationID string) error

	//获取hr的所有面试
	GetInterviewsByUserID(ctx context.Context, userID string) ([]entity.Interview, error)

	//获取某面试的聊天记录
	GetChatRecordByConversationID(ctx context.Context, conversationID string) ([]entity.Message, error)

	//保存整个聊天记录
	SaveAllMessages(ctx context.Context, conversationID string, messages []entity.Message) error

	//在末尾添加单条聊天记录
	AddMessage(ctx context.Context, conversationID string, message entity.Message) error

	//删除某面试
	DeleteInterview(ctx context.Context, conversationID string) error
}
