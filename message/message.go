package message

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

// Messenger :
type Messenger interface {
	GetReply(
		context.Context,
		*openai.Client,
		string,
	) (*openai.ChatCompletionMessage, error)
}

func createMessages(
	message string,
	history []openai.ChatCompletionMessage,
) []openai.ChatCompletionMessage {
	return append(history, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: message,
	})
}

func getReply(
	ctx context.Context,
	client *openai.Client,
	messages []openai.ChatCompletionMessage,
) (*openai.ChatCompletionMessage, error) {
	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: messages,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("createChatCompletion: %v", err)
	}

	expectedReplyLength := 1
	if len(resp.Choices) != expectedReplyLength {
		return nil, fmt.Errorf("length of choices should be %d", expectedReplyLength)
	}

	return &resp.Choices[expectedReplyLength-1].Message, nil
}
