package openai

import (
	"errors"

	"github.com/sashabaranov/go-openai"
)

var ErrNoChat = errors.New("no chat found")

type Chats []*Chat

func (c Chats) Create(id int64, msg openai.ChatCompletionMessage) (Chats, *Chat) {
	chat := &Chat{
		Id:       id,
		Messages: []openai.ChatCompletionMessage{msg},
	}

	return append(c, chat), chat
}

func (c Chats) Find(id int64) (*Chat, error) {
	for _, chat := range c {
		if chat.Id == id {
			return chat, nil
		}
	}

	return nil, ErrNoChat
}

func (c Chats) Reset(id int64) {
	for _, chat := range c {
		if chat.Id == id {
			chat.Messages = []openai.ChatCompletionMessage{}
		}
	}
}

type Chat struct {
	Id       int64
	Messages []openai.ChatCompletionMessage
}

func (c *Chat) AddMessage(msg openai.ChatCompletionMessage) {
	c.Messages = append(c.Messages, msg)
}
