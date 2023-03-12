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

		var messenger message.Messenger
		{
			messagesPath, acceptMessageHistory := os.LookupEnv("CHATTY_MESSAGES_PATH")
			if acceptMessageHistory {
				file, err := os.OpenFile(messagesPath, os.O_RDWR|os.O_CREATE, 0666)
				if err != nil {
					return fmt.Errorf("create or open file: %w", err)
				}
				defer file.Close()
				messenger = message.NewHistoricalMessenger(file)
			} else {
				messenger = message.NewIndependentMessenger()
			}
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

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
