package message

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/sashabaranov/go-openai"
)

// NewHistoricalMessenger :
func NewHistoricalMessenger(store HistoryStore) Messenger {
	return &historicalMessenger{
		store: store,
	}
}

// historicalMessenger :
type historicalMessenger struct {
	store HistoryStore
}

func (m *historicalMessenger) GetReply(
	ctx context.Context,
	client *openai.Client,
	message string,
) (*openai.ChatCompletionMessage, error) {
	messageHistory, err := m.load()
	if err != nil {
		return nil, fmt.Errorf("load message history: %w", err)
	}

	messages := createMessages(message, messageHistory)

	reply, err := getReply(ctx, client, messages)
	if err != nil {
		return nil, fmt.Errorf("get reply: %w", err)
	}

	if err := m.update(append(messages, *reply)); err != nil {
		return nil, fmt.Errorf("update message history: %w", err)
	}

	return reply, nil
}

// HistoryStore :
type HistoryStore interface {
	io.Reader
	io.WriterAt
}

// load :
func (m *historicalMessenger) load() ([]openai.ChatCompletionMessage, error) {
	var history []openai.ChatCompletionMessage
	if err := json.NewDecoder(m.store).Decode(&history); err != nil && !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return history, nil
}

// update :
func (m *historicalMessenger) update(messageHistory []openai.ChatCompletionMessage) error {
	bytes, err := json.MarshalIndent(messageHistory, "", " ")
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}
	if _, err := m.store.WriteAt(bytes, 0); err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	return nil
}
