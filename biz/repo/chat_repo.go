package repo

import (
	"ai_interview/biz/entity"
	"context"
	"gorm.io/datatypes"
)

type ChatRepo interface {

	//新的关于聊天记录的方法

	//新建面试
	CreateInterview(ctx context.Context, userID, talentID, roomID string) error

	//结束面试
	EndInterview(ctx context.Context, roomID, userID string) error

	//获取hr的所有面试
	GetInterviewsByUserID(ctx context.Context, userID string) ([]entity.Interview, error)

	//保存面试聊天记录
	SaveInterviewMessage(ctx context.Context, message entity.InterviewMessage) error

	//获取某面试的聊天记录
	GetInterviewMessageByRoomID(ctx context.Context, roomID string) ([]entity.InterviewMessage, error)

	//给某条聊天记录加上标签
	AddTagToMessage(ctx context.Context, messageID string, tag string, roomID string) error

	//检查面试房间是否存在
	CheckRoom(ctx context.Context, roomID, userID string) (bool, error)

	//根据roomID获取面试
	GetInterviewByRoomID(ctx context.Context, roomID, userID string) (*entity.Interview, error)

	//根据TalentID来获取面试
	GetInterviewByTalentID(ctx context.Context, talentID, userID string) (*entity.Interview, error)

	//保存生成的人才报告
	SaveTalentReport(ctx context.Context, report datatypes.JSON, userID, talentID string) error
}
