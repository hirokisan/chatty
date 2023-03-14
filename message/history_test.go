package message

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/hirokisan/chatty/testhelper"
	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testWriter struct {
	bytes.Buffer
}

func (w *testWriter) WriteAt(b []byte, _ int64) (n int, err error) {
	w.Buffer = *bytes.NewBuffer(b)
	return len(b), nil
}

func TestHistoricalMessenger_GetReply(t *testing.T) {
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

	var store testWriter
	messenger := NewHistoricalMessenger(&store)

	message := "my name is chatty"
	got, err := messenger.GetReply(ctx, client, message)
	require.NoError(t, err)

	{
		assert.Equal(t, respBody.Choices[0].Message, *got)

		var messageHistory []openai.ChatCompletionMessage
		require.NoError(t, json.NewDecoder(&store).Decode(&messageHistory))

		assert.Equal(t, message, messageHistory[0].Content)
	}
}
