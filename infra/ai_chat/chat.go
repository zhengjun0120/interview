package ai_chat

import (
	"ai_interview/biz/entity"
	"ai_interview/conf"
	"ai_interview/pkg/zlog"
	"context"
	"fmt"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model/responses"
)

func (c *ChatClient) Chat(ctx context.Context, messages []entity.Message) (string, error) {

	reqMessage := Message2ReqMessage(messages)

	req := model.CreateChatCompletionRequest{
		Model:    conf.GetConfig().AiChat.ModelId,
		Messages: reqMessage,
		Thinking: &model.Thinking{
			Type: model.ThinkingTypeAuto,
		},
	}

	resp, err := c.Client.CreateChatCompletion(ctx, req)

	if err != nil {
		return "", err
	}
	return *resp.Choices[0].Message.Content.StringValue, nil

}

// 解析文档输出优势等内容
func (c *ChatClient) DocxToTalentDataChat(ctx context.Context, docxUrl, jobProfile string) (string, error) {
	systemMessage := &responses.ItemInputMessage{
		Role: responses.MessageRole_system,
		Content: []*responses.ContentItem{
			{
				Union: &responses.ContentItem_Text{
					Text: &responses.ContentItemText{
						Type: responses.ContentItemType_input_text,
						Text: conf.GetConfig().AiChat.ResumeAnalysisSystemPrompt,
					},
				},
			},
		},
	}

	inputMessage := &responses.ItemInputMessage{
		Role: responses.MessageRole_user,
	}

	if jobProfile != "" {
		content := []*responses.ContentItem{
			{
				Union: &responses.ContentItem_File{
					File: &responses.ContentItemFile{
						Type:    responses.ContentItemType_input_file,
						FileUrl: &docxUrl,
					},
				},
			},
			{
				Union: &responses.ContentItem_Text{
					Text: &responses.ContentItemText{
						Type: responses.ContentItemType_input_text,
						Text: "所有岗位画像如下\n\n" + jobProfile,
					},
				},
			},
		}
		inputMessage.Content = content
	} else {
		Content := []*responses.ContentItem{
			{
				Union: &responses.ContentItem_File{
					File: &responses.ContentItemFile{
						Type:    responses.ContentItemType_input_file,
						FileUrl: &docxUrl,
					},
				},
			},
		}
		inputMessage.Content = Content
	}

	resp, err := c.Client.CreateResponses(ctx, &responses.ResponsesRequest{
		Model: conf.GetConfig().AiChat.ModelId,
		Input: &responses.ResponsesInput{
			Union: &responses.ResponsesInput_ListValue{
				ListValue: &responses.InputItemList{
					ListValue: []*responses.InputItem{
						{
							Union: &responses.InputItem_InputMessage{
								InputMessage: systemMessage,
							},
						},
						{
							Union: &responses.InputItem_InputMessage{
								InputMessage: inputMessage,
							},
						},
					},
				},
			},
		},
		Thinking: &responses.ResponsesThinking{Type: responses.ThinkingType_enabled.Enum()},
	})

	if err != nil {
		return "", err
	}

	for _, output := range resp.Output {
		switch v := output.Union.(type) {
		case *responses.OutputItem_OutputMessage:
			outMessage, ok := v.OutputMessage.Content[0].Union.(*responses.OutputContentItem_Text)
			if !ok {
				zlog.Errorf("docxChat 内容输出错误")
				fmt.Println(v)
				return "", fmt.Errorf("docxChat 内容输出错误")
			}

			return outMessage.Text.Text, nil // 返回输出消息的文本内容
		}
	}

	return "", fmt.Errorf("docxChat 无输出")
}

func (c *ChatClient) DocxToResumeStrChat(ctx context.Context, docxUrl string) (string, error) {
	systemMessage := &responses.ItemInputMessage{
		Role: responses.MessageRole_system,
		Content: []*responses.ContentItem{
			{
				Union: &responses.ContentItem_Text{
					Text: &responses.ContentItemText{
						Type: responses.ContentItemType_input_text,
						Text: conf.GetConfig().AiChat.ResumeUrlToStrPrompt,
					},
				},
			},
		},
	}

	inputMessage := &responses.ItemInputMessage{
		Role: responses.MessageRole_user,
		Content: []*responses.ContentItem{
			{
				Union: &responses.ContentItem_File{
					File: &responses.ContentItemFile{
						Type:    responses.ContentItemType_input_file,
						FileUrl: &docxUrl,
					},
				},
			},
			//{
			//	Union: &responses.ContentItem_Text{
			//		Text: &responses.ContentItemText{
			//			Type: responses.ContentItemType_input_text,
			//			Text: "请将文档内容转换为纯文本，并返回纯文本内容。",
			//		},
			//	}
			//},
		},
	}

	// 调用CreateResponses方法
	resp, err := c.Client.CreateResponses(ctx, &responses.ResponsesRequest{
		// 模型ID
		Model: conf.GetConfig().AiChat.ModelId,
		// 输入使用&responses.ResponsesInput
		Input: &responses.ResponsesInput{
			// 使用&responses.ResponsesInput_ListValue
			Union: &responses.ResponsesInput_ListValue{
				// 输入项列表
				ListValue: &responses.InputItemList{
					// 输入项
					ListValue: []*responses.InputItem{
						{
							Union: &responses.InputItem_InputMessage{
								InputMessage: systemMessage,
							},
						},
						{
							Union: &responses.InputItem_InputMessage{
								InputMessage: inputMessage,
							},
						},
					},
				},
			},
		},
		Thinking: &responses.ResponsesThinking{Type: responses.ThinkingType_disabled.Enum()},
	})

	if err != nil {
		return "", err
	}

	for _, output := range resp.Output {
		switch v := output.Union.(type) {
		case *responses.OutputItem_OutputMessage:
			outMessage, ok := v.OutputMessage.Content[0].Union.(*responses.OutputContentItem_Text)
			if !ok {
				zlog.Errorf("docxChat 内容输出错误")
				fmt.Println(v)
				return "", fmt.Errorf("docxChat 内容输出错误")
			}

			return outMessage.Text.Text, nil // 返回输出消息的文本内容
		}
	}

	return "", fmt.Errorf("docxChat 无输出")
}

// 获取面试AI建议
func (c *ChatClient) InterviewAiSuggestionChat(ctx context.Context, interviewMessage []entity.InterviewMessage, resumeStr, resumeUrl string) (string, error) {
	inputMessageStr, err := buildInterviewMessage(interviewMessage, resumeStr)
	if err != nil {
		return "", err
	}

	// 如果 resumeStr 不为空，则将 resumeUrl 置为空 使用文本而不是连接
	if resumeStr != "" {
		resumeUrl = ""
	}
	resp, err := c.buildInputChat(ctx, conf.GetConfig().AiChat.InterviewAISuggestionPrompt, inputMessageStr, resumeUrl)
	if err != nil {
		return "", err
	}
	return resp, nil
}

// 生成面试报告
func (c *ChatClient) InterviewMessageAnalyseChat(ctx context.Context, interviewMessage []entity.InterviewMessage) (string, error) {

	inputMessageStr, err := buildInterviewMessage(interviewMessage, "")
	if err != nil {
		return "", err
	}

	resp, err := c.buildInputChat(ctx, conf.GetConfig().AiChat.InterviewMessageAnalysePrompt, inputMessageStr, "")
	if err != nil {
		return "", err
	}
	return resp, nil
}

// 统一构建输入消息 并获取输出 适应多模态
func (c *ChatClient) buildInputChat(ctx context.Context, systemMessageStr string, userMessageStr string, fileUrl string) (string, error) {
	// 构建上下文，使用&responses.ItemInputMessage
	systemMessage := &responses.ItemInputMessage{
		// 角色：系统
		Role: responses.MessageRole_system,
		// 内容：系统消息 构建内容使用 &responses.ContentItem
		Content: []*responses.ContentItem{
			{
				// 内容类型：文本
				Union: &responses.ContentItem_Text{
					Text: &responses.ContentItemText{
						Type: responses.ContentItemType_input_text,
						Text: systemMessageStr,
					},
				},
			},
		},
	}

	userInputMessage := &responses.ItemInputMessage{
		Role: responses.MessageRole_user,
	}

	// 构建用户输入消息
	var content []*responses.ContentItem
	// 根据文件URL是否为空，构建不同内容 多模态判断
	if fileUrl == "" {
		content = []*responses.ContentItem{
			{
				Union: &responses.ContentItem_Text{
					Text: &responses.ContentItemText{
						Type: responses.ContentItemType_input_text,
						Text: userMessageStr,
					},
				},
			},
		}
	} else {
		content = []*responses.ContentItem{
			{
				Union: &responses.ContentItem_File{
					File: &responses.ContentItemFile{
						Type:    responses.ContentItemType_input_file,
						FileUrl: &fileUrl,
					},
				},
			},
			{
				Union: &responses.ContentItem_Text{
					Text: &responses.ContentItemText{
						Type: responses.ContentItemType_input_text,
						Text: userMessageStr,
					},
				},
			},
		}
	}

	userInputMessage.Content = content

	resp, err := c.Client.CreateResponses(ctx, &responses.ResponsesRequest{
		Model: conf.GetConfig().AiChat.ModelId,
		Input: &responses.ResponsesInput{
			Union: &responses.ResponsesInput_ListValue{
				ListValue: &responses.InputItemList{
					ListValue: []*responses.InputItem{
						{
							Union: &responses.InputItem_InputMessage{
								InputMessage: systemMessage,
							},
						},
						{
							Union: &responses.InputItem_InputMessage{
								InputMessage: userInputMessage,
							},
						},
					},
				},
			},
		},
		// 思考过程：禁用
		Thinking: &responses.ResponsesThinking{Type: responses.ThinkingType_disabled.Enum()},
	})

	if err != nil {
		return "", err
	}

	// 遍历输出项，使用 switch 语句判断输出项类型
	for _, output := range resp.Output {
		switch v := output.Union.(type) {
		// 输出项类型：输出消息 输出消息内容使用 &responses.OutputContentItem_Text
		case *responses.OutputItem_OutputMessage:
			outMessage, ok := v.OutputMessage.Content[0].Union.(*responses.OutputContentItem_Text)
			if !ok {
				zlog.Errorf("BuildInputChat 内容输出错误")
				fmt.Println(v)
				return "", fmt.Errorf("BuildInputChat 内容输出错误")
			}

			return outMessage.Text.Text, nil // 返回输出消息的文本内容
		}
	}
	return "", fmt.Errorf("BuildInputChat 无输出")
}
