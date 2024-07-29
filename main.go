package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

const (
	RemindCommand = "!r"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(GetEnv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}
	u := tgbotapi.NewUpdate(0)
	// debug logs
	// bot.Debug = true

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatalf("Error getting updates: %v", err)
	}

	for update := range updates {
		isValidCommand := strings.HasPrefix(update.Message.Text, RemindCommand)
		if !isValidCommand {
			fmt.Println("Error. Invalid message")
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "The only command to use for now is '!r' to set a reminder. The message you sent is not valid.")
			_, err := bot.Send(msg)
			if err != nil {
				log.Printf("Error sending message: %v", err)
			}
		} else {
			index := strings.Index(update.Message.Text, RemindCommand)
			if index == -1 {
				fmt.Println("Command not found in text.")
				return
			}

			textMessage := update.Message.Text[index+len(RemindCommand):]
			textMessage = strings.TrimSpace(textMessage)
			fmt.Println("New message: " + textMessage)

			if update.Message == nil { // Ignore any non-Message updates.
				continue
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "New message: "+textMessage)
			_, err := bot.Send(msg)
			if err != nil {
				log.Printf("Error sending message: %v", err)
			}
		}

	}

}

func GetEnv(key string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	token := os.Getenv(key)
	if token == "" {
		log.Fatal(key + " is not set in the environment")
	}

	return token
}
