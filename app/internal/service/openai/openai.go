package openai

import (
	"context"
	"errors"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

type Service struct {
	client *openai.Client

	chats Chats
}

func NewService(cfg *Config) (*Service, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}

	client := openai.NewClient(cfg.Token)

	return &Service{client: client, chats: []*Chat{}}, nil
}

func (s *Service) NewConversation(_ context.Context, id int64) {
	s.chats.Reset(id)
}

func (s *Service) ChatCompletion(ctx context.Context, id int64, prompt string) (string, error) {
	chat, err := s.chats.Find(id)
	if err != nil && !errors.Is(err, ErrNoChat) {
		return "", err
	}

	newMessage := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: prompt,
	}

	if errors.Is(err, ErrNoChat) {
		s.chats, chat = s.chats.Create(id, newMessage)
	} else {
		chat.AddMessage(newMessage)
	}

	resp, err := s.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    openai.GPT4o,
			Messages: chat.Messages,
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func (s *Service) GenerateImagePrompt(ctx context.Context, prompt string) (string, error) {
	prompt = fmt.Sprintf("this is human description about image what does he want, "+
		"now generate prompt for neural network, to generate image. use a lot of words. "+
		"be precise and masterpiece! here is user prompt: %s", prompt)
	newMessage := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: prompt,
	}

	resp, err := s.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    openai.GPT4o,
			Messages: []openai.ChatCompletionMessage{newMessage},
		},
	)
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
