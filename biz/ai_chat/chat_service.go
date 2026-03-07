package ai_chat

import (
	"ai_interview/biz/entity"
	"context"
)

type IChatService interface {
	// 基础聊天已弃用
	Chat(ctx context.Context, messages []entity.Message) (string, error)
	// Docx转人才数据聊天
	DocxToTalentDataChat(ctx context.Context, dockUrl, jobProfile string) (string, error)
	// Docx转简历字符串聊天
	DocxToResumeStrChat(ctx context.Context, dockUrl string) (string, error)
	// 面试AI建议聊天
	InterviewAiSuggestionChat(ctx context.Context, interviewMessage []entity.InterviewMessage, resumeStr, resumeUrl string) (string, error)
	// 生成人才报告
	InterviewMessageAnalyseChat(ctx context.Context, interviewMessage []entity.InterviewMessage) (string, error)
}
