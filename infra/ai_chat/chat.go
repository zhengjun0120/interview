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
