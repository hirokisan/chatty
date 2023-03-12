package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hirokisan/chatty/message"
	openai "github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

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

		messenger, err := getMessenger()
		if err != nil {
			return fmt.Errorf("get messenger: %w", err)
		}

		msg := strings.Join(args, " ")
		reply, err := messenger.GetReply(ctx, client, msg)
		if err != nil {
			return fmt.Errorf("get reply: %w", err)
		}

		fmt.Println(reply.Content)
		return nil
	},
}

func getMessenger() (message.Messenger, error) {
	messagesPath, acceptMessageHistory := os.LookupEnv("CHATTY_MESSAGES_PATH")
	if acceptMessageHistory {
		file, err := os.OpenFile(messagesPath, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			return nil, fmt.Errorf("create or open file: %w", err)
		}
		defer file.Close()

		return message.NewHistoricalMessenger(file), nil
	}
	return message.NewIndependentMessenger(), nil
}
