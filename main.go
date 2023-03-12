package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	openai "github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "chatty",
	Short: "Would you like to have a little chat with us between jobs?",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		key, ok := os.LookupEnv("OPEN_AI_KEY")
		if !ok {
			return fmt.Errorf("OPEN_AI_KEY should be set as environmental variable")
		}
		client := openai.NewClient(key)

		messagesPath, acceptMessageHistory := os.LookupEnv("CHATTY_MESSAGES_PATH")

		var messageHistory []openai.ChatCompletionMessage
		if acceptMessageHistory {
			if err := loadMessageHistory(messagesPath, &messageHistory); err != nil {
				return fmt.Errorf("load message history: %w", err)
			}
		}

		message := strings.Join(args, " ")
		messages := createMessages(message, messageHistory)

		reply, err := getReply(ctx, client, messages)
		if err != nil {
			return fmt.Errorf("get reply: %w", err)
		}

		if acceptMessageHistory {
			if err := updateMessageHistory(messagesPath, append(messages, *reply)); err != nil {
				return fmt.Errorf("update message history: %w", err)
			}
		}

		fmt.Println(reply.Content)
		return nil
	},
}

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func loadMessageHistory(
	path string,
	messageHistory *[]openai.ChatCompletionMessage,
) error {
	file, err := os.OpenFile(path, os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("create or open file: %w", err)
	}
	defer file.Close()
	if err := json.NewDecoder(file).Decode(&messageHistory); err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("decode: %w", err)
	}
	return nil
}

func updateMessageHistory(
	path string,
	messageHistory []openai.ChatCompletionMessage,
) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	bytes, err := json.MarshalIndent(messageHistory, "", " ")
	if _, err := file.Write(bytes); err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	return nil
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
