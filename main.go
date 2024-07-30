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
	StartCommand  = "/start"
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
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatalf("Error getting updates: %v", err)
	}

	redisInstance, err := StartRedis(dbInstance, bot)
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
		return
	}

	for update := range updates {
		isValidCommand := strings.HasPrefix(update.Message.Text, RemindCommand)
		if !isValidCommand {
			fmt.Println("Error. Invalid message")
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "The only command to use for now is '/r' to set a reminder. The message you sent is not valid.")
			_, err := bot.Send(msg)
			if err != nil {
				panic("Error trying to send message: " + err.Error())
			}
		} else {
			title, notifyTime, duration, err := ParseMessage(update.Message.Text)
			if err != nil {
				panic("Error trying to parse message: " + err.Error())
			}

			if update.Message == nil { // Ignore any non-Message updates.
				continue
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Reminder setted: '"+title+"'")
			result, err := dbInstance.InsertNewReminder(title, notifyTime, update.Message.Chat.ID)
			if err != nil {
				panic("Error trying to insert new reminder in db: " + err.Error())
			}

			err = redisInstance.Set(strconv.Itoa(result.ID), "", duration)
			if err != nil {
				panic("Error setting value in Redis:" + err.Error())
			}

			_, err = bot.Send(msg)
			if err != nil {
				panic("Error trying to send message: " + err.Error())
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

func ParseMessage(message string) (text string, notifyTime time.Time, duration time.Duration, err error) {
	re := regexp.MustCompile(`/r\s+"(.*?)"\s+in\s+"(.*?)"`)
	matches := re.FindStringSubmatch(message)

	if len(matches) < 3 {
		fmt.Println("Matches:", matches)
		return "", time.Time{}, 0, fmt.Errorf("invalid message format")
	}

	text = matches[1]
	durationStr := matches[2]
	durationStr = strings.TrimSpace(durationStr)

	if strings.HasSuffix(durationStr, "h") {
		hours, err := strconv.Atoi(strings.TrimSuffix(durationStr, "h"))
		if err != nil {
			return "", time.Time{}, 0, fmt.Errorf("error parsing hours: %w", err)
		}
		duration = time.Duration(hours) * time.Hour
	} else if strings.HasSuffix(durationStr, "m") {
		minutes, err := strconv.Atoi(strings.TrimSuffix(durationStr, "m"))
		if err != nil {
			return "", time.Time{}, 0, fmt.Errorf("error parsing minutes: %w", err)
		}
		duration = time.Duration(minutes) * time.Minute
	} else if strings.HasSuffix(durationStr, "s") {
		seconds, err := strconv.Atoi(strings.TrimSuffix(durationStr, "s"))
		if err != nil {
			return "", time.Time{}, 0, fmt.Errorf("error parsing seconds: %w", err)
		}
		duration = time.Duration(seconds) * time.Second
	} else {
		return "", time.Time{}, 0, fmt.Errorf("invalid duration format")
	}

	notifyTime = time.Now().Add(duration)

	index := strings.Index(text, RemindCommand)
	if index == -1 {
		fmt.Println("Command not found in text.")
		return
	}
	textMessage := text[index+len(RemindCommand):]
	textMessage = strings.TrimSpace(textMessage)

	return textMessage, notifyTime, duration, nil
}

func SendTelegramMessage(bot *tgbotapi.BotAPI, chatID int64, message string) error {
	msg := tgbotapi.NewMessage(chatID, message)
	_, err := bot.Send(msg)
	if err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}
	return nil
}
