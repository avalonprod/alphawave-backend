package openai

import (
	"context"

	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/types"
	openai "github.com/sashabaranov/go-openai"
)

const GPT_MODEL = "gpt-3.5-turbo"

type OpenAiAPI struct {
	token string
	url   string
}

func NewOpenAiAPI(token string, url string) *OpenAiAPI {
	return &OpenAiAPI{
		token: token,
		url:   url,
	}
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
type ResponseData struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	// Model   string             `json:"model"`
	Choices []messagesResponse `json:"choices"`
	Usage   Usage              `json:"usage"`
}
type Response struct {
	Data ResponseData `json:"data"`
}

type messagesResponse struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type OutputMessage struct {
	Message Message
}

// func (o *OpenAiAPI) NewMessage(messages []types.Message) (OutputMessage, error) {
// 	client := openai.NewClient(o.token)

// 	var reqMessages []openai.ChatCompletionMessage
// 	reqMessages = append(reqMessages, openai.ChatCompletionMessage{
// 		Role:    "system",
// 		Content: "When someone says hello to you, you should say hello too. When they beat you up, tell them you're AlphaWave INC's artificial intelligence assistant.\nAnd ask me how I can help.",
// 	})

// 	for _, item := range messages {
// 		reqMessages = append(reqMessages, openai.ChatCompletionMessage{
// 			Role:    item.Role,
// 			Content: item.Content,
// 		})
// 	}

// 	resp, err := client.CreateChatCompletion(
// 		context.Background(),
// 		openai.ChatCompletionRequest{
// 			Model:    openai.GPT3Dot5Turbo,
// 			Messages: reqMessages,
// 		},
// 	)

// 	if err != nil {
// 		return OutputMessage{}, err
// 	}

// 	return OutputMessage{
// 		Message: Message{
// 			Role:    resp.Choices[0].Message.Role,
// 			Content: resp.Choices[0].Message.Content,
// 		},
// 	}, nil
// }

func (o *OpenAiAPI) NewMessage(messages []types.Message) (*types.MessageOutput, error) {
	client := openai.NewClient(o.token)

	var reqMessages []openai.ChatCompletionMessage
	reqMessages = append(reqMessages, openai.ChatCompletionMessage{
		Role:    "system",
		Content: "When someone says hello to you, you should say hello too. When they beat you up, tell them you're AlphaWave INC's artificial intelligence assistant.\nAnd ask me how I can help.",
	})

	for _, item := range messages {
		reqMessages = append(reqMessages, openai.ChatCompletionMessage{
			Role:    item.Role,
			Content: item.Content,
		})
	}

	req := openai.ChatCompletionRequest{
		Model:    openai.GPT4,
		Messages: reqMessages,
		Stream:   true,
	}

	stream, err := client.CreateChatCompletionStream(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return &types.MessageOutput{
		Role:   "assistent",
		Stream: *stream,
	}, nil

}
