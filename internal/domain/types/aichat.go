package types

import "github.com/sashabaranov/go-openai"

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type MessageOutput struct {
	Role   string
	Stream openai.ChatCompletionStream
}
