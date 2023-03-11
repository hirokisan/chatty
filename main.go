package main

import (
	"context"
	"fmt"
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

		message := strings.Join(args, " ")

		reply, err := getReply(ctx, client, message)
		if err != nil {
			return fmt.Errorf("get reply: %w", err)
		}

		fmt.Println(reply)
		return nil
	},
}

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func getReply(
	ctx context.Context,
	client *openai.Client,
	message string,
) (string, error) {
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: message,
		},
	}
	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: messages,
		},
	)
	if err != nil {
		return "", fmt.Errorf("createChatCompletion: %v", err)
	}

	if len(resp.Choices) < len(messages) {
		return "", fmt.Errorf("length of choices should be more than %d but got %d", len(messages), len(resp.Choices))
	}

	reply := strings.ReplaceAll(resp.Choices[len(messages)-1].Message.Content, "\n", "")

	return reply, nil
}
