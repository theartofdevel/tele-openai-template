package openai

import (
	"testing"

	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
)

func TestChats_Create(t *testing.T) {
	chats := Chats{}
	msg := openai.ChatCompletionMessage{
		Role:    "user",
		Content: "Hello",
	}

	newChats, chat := chats.Create(1, msg)

	assert.Len(t, newChats, 1)
	assert.Equal(t, int64(1), chat.Id)
	assert.Len(t, chat.Messages, 1)
	assert.Equal(t, msg, chat.Messages[0])
}

func TestChats_Find(t *testing.T) {
	msg := openai.ChatCompletionMessage{
		Role:    "user",
		Content: "Hello",
	}
	chats := Chats{
		{
			Id:       1,
			Messages: []openai.ChatCompletionMessage{msg},
		},
	}

	tests := []struct {
		name    string
		id      int64
		want    *Chat
		wantErr error
	}{
		{
			name: "existing chat",
			id:   1,
			want: &Chat{
				Id:       1,
				Messages: []openai.ChatCompletionMessage{msg},
			},
			wantErr: nil,
		},
		{
			name:    "non-existing chat",
			id:      2,
			want:    nil,
			wantErr: ErrNoChat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := chats.Find(tt.id)
			assert.Equal(t, tt.wantErr, err)
			if tt.want != nil {
				assert.Equal(t, tt.want.Id, got.Id)
				assert.Equal(t, tt.want.Messages, got.Messages)
			} else {
				assert.Nil(t, got)
			}
		})
	}
}

func TestChats_Reset(t *testing.T) {
	msg := openai.ChatCompletionMessage{
		Role:    "user",
		Content: "Hello",
	}
	chats := Chats{
		{
			Id:       1,
			Messages: []openai.ChatCompletionMessage{msg},
		},
		{
			Id:       2,
			Messages: []openai.ChatCompletionMessage{msg},
		},
	}

	chats.Reset(1)

	chat1, _ := chats.Find(1)
	assert.Len(t, chat1.Messages, 0)

	chat2, _ := chats.Find(2)
	assert.Len(t, chat2.Messages, 1)
	assert.Equal(t, msg, chat2.Messages[0])
}

func TestChat_AddMessage(t *testing.T) {
	chat := &Chat{
		Id:       1,
		Messages: []openai.ChatCompletionMessage{},
	}

	msg1 := openai.ChatCompletionMessage{
		Role:    "user",
		Content: "Hello",
	}
	msg2 := openai.ChatCompletionMessage{
		Role:    "assistant",
		Content: "Hi there",
	}

	chat.AddMessage(msg1)
	assert.Len(t, chat.Messages, 1)
	assert.Equal(t, msg1, chat.Messages[0])

	chat.AddMessage(msg2)
	assert.Len(t, chat.Messages, 2)
	assert.Equal(t, msg2, chat.Messages[1])
}
