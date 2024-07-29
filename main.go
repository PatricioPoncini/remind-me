package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

const (
	RemindCommand = "/r"
)

func main() {
	dbInstance, err := NewDB(GetEnv("DB_KEY"))
	if err != nil {
		log.Fatalf("Error creating database instance: %v", err)
		return
	}
	dbInstance.checkInitialConditions()

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
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "The only command to use for now is '/r' to set a reminder. The message you sent is not valid.") // bold message --> '/r'
			_, err := bot.Send(msg)
			if err != nil {
				log.Printf("Error sending message: %v", err)
			}
		} else {
			fmt.Println("New message: " + update.Message.Text)
			title, notifyTime, err := ParseMessage(update.Message.Text)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			if update.Message == nil { // Ignore any non-Message updates.
				continue
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "New message: "+title)
			err = dbInstance.InsertNewReminder(title, notifyTime)
			if err != nil {
				panic(err)
			}
			_, err = bot.Send(msg)
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

func ParseMessage(message string) (text string, notifyTime time.Time, err error) {
	re := regexp.MustCompile(`/r\s+"(.*?)"\s+en\s+"(.*?)"`)
	matches := re.FindStringSubmatch(message)

	if len(matches) < 3 {
		fmt.Println("Matches:", matches) // Debug print
		return "", time.Time{}, fmt.Errorf("invalid message format")
	}

	text = matches[1]

	durationStr := matches[2]
	durationStr = strings.TrimSpace(durationStr)

	var duration time.Duration
	if strings.HasSuffix(durationStr, "h") {
		hours, err := strconv.Atoi(strings.TrimSuffix(durationStr, "h"))
		if err != nil {
			return "", time.Time{}, fmt.Errorf("error parsing hours: %w", err)
		}
		duration = time.Duration(hours) * time.Hour
	} else if strings.HasSuffix(durationStr, "m") {
		minutes, err := strconv.Atoi(strings.TrimSuffix(durationStr, "m"))
		if err != nil {
			return "", time.Time{}, fmt.Errorf("error parsing minutes: %w", err)
		}
		duration = time.Duration(minutes) * time.Minute
		// TODO: parse to seconds
	} else {
		return "", time.Time{}, fmt.Errorf("invalid duration format")
	}

	notifyTime = time.Now().Add(duration)

	index := strings.Index(text, RemindCommand)
	if index == -1 {
		fmt.Println("Command not found in text.")
		return
	}
	textMessage := text[index+len(RemindCommand):]
	textMessage = strings.TrimSpace(textMessage)

	return textMessage, notifyTime, nil
}

// TODO: FIND REDIS KEY BY MESSAGE ID
/*result, err := dbInstance.GetReminderById(1)
if err != nil {
	fmt.Println("Error:", err)
	return
}*/
