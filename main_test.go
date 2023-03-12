package main

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/hirokisan/chatty/testhelper"
	openai "github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetReply(t *testing.T) {
	ctx := context.Background()

	path := "/chat/completions"
	method := http.MethodPost
	status := http.StatusOK
	respBody := openai.ChatCompletionResponse{
		ID:      "chatcmpl-6t43kFP87uHaw2DGOlnoVF7806A90",
		Object:  "chat.completion",
		Created: 1678581664,
		Model:   "gpt-3.5-turbo-0301",
		Choices: []openai.ChatCompletionChoice{
			{
				Index: 0,
				Message: openai.ChatCompletionMessage{
					Role:    "assistant",
					Content: "Hello, Chatty! How can I assist you today?",
				},
				FinishReason: "",
			},
		},
		Usage: openai.Usage{
			PromptTokens:     16,
			CompletionTokens: 64,
			TotalTokens:      80,
		},
	}
	bytesBody, err := json.Marshal(respBody)
	require.NoError(t, err)

	server, teardown := testhelper.NewServer(
		testhelper.WithHandlerOption(path, method, status, bytesBody),
	)
	defer teardown()

	client := testhelper.NewTestClient(server.URL)

	message := "my name is chatty"
	got, err := getReply(ctx, client, message)
	require.NoError(t, err)

	assert.Equal(t, respBody.Choices[0].Message.Content, got)
}
