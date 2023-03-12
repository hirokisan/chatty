package message

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

type independentMessenger struct{}

// NewIndependentMessenger :
func NewIndependentMessenger() Messenger {
	return &independentMessenger{}
}

func (m *independentMessenger) GetReply(
	ctx context.Context,
	client *openai.Client,
	message string,
) (*openai.ChatCompletionMessage, error) {
	messages := createMessages(message, nil)

	reply, err := getReply(ctx, client, messages)
	if err != nil {
		return nil, fmt.Errorf("get reply: %w", err)
	}

	return reply, nil
}
