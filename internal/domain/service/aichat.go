package service

import (
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/types"
)

type openAI interface {
	NewMessage(messages []types.Message) (*types.MessageOutput, error)
}

type AiChatService struct {
	openAI openAI
}

func NewAiChatService(openAI openAI) *AiChatService {
	return &AiChatService{
		openAI: openAI,
	}
}

func (s *AiChatService) NewMessage(messages []types.Message) (*types.MessageOutput, error) {

	message, err := s.openAI.NewMessage(messages)

	if err != nil {
		return &types.MessageOutput{}, err
	}

	return message, nil
}
