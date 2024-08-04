package bot

import (
	"fmt"
	"log"
	"remind_me/src/db"
	"remind_me/src/redis"
	"remind_me/src/utils"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type TelegramBot struct {
	instance *tgbotapi.BotAPI
}

func NewTelegramBot() (*TelegramBot, error) {
	bot, err := tgbotapi.NewBotAPI(utils.GetEnv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		return nil, fmt.Errorf("error creating bot: %v", err)
	}
	return &TelegramBot{instance: bot}, nil
}

func (tb *TelegramBot) Start(dbInstance *db.DB, redisInstance *redis.Redis) {
	u := tgbotapi.NewUpdate(0)
	updates, err := tb.instance.GetUpdatesChan(u)
	if err != nil {
		log.Fatalf("Error getting updates: %v", err)
	}

	for update := range updates {
		isValidCommand := strings.HasPrefix(update.Message.Text, utils.RemindCommand)
		if !isValidCommand {
			fmt.Println("Error. Invalid message")
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "The only command to use for now is '/r' to set a reminder. The message you sent is not valid.")
			_, err := tb.instance.Send(msg)
			if err != nil {
				panic("Error trying to send message: " + err.Error())
			}
		} else {
			title, notifyTime, duration, err := utils.ParseMessage(update.Message.Text)
			if err != nil {
				panic("Error trying to parse message: " + err.Error())
			}

			if update.Message == nil { // Ignore any non-Message updates.
				continue
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Reminder set: '"+title+"'")
			result, err := dbInstance.InsertNewReminder(title, notifyTime, update.Message.Chat.ID)
			if err != nil {
				panic("Error trying to insert new reminder in db: " + err.Error())
			}

			err = redisInstance.Set(strconv.Itoa(result.ID), "", duration)
			if err != nil {
				panic("Error setting value in Redis:" + err.Error())
			}

			_, err = tb.instance.Send(msg)
			if err != nil {
				panic("Error trying to send message: " + err.Error())
			}
		}
	}
}

func (tb *TelegramBot) SendTelegramMessage(chatID int64, message string) error {
	msg := tgbotapi.NewMessage(chatID, message)
	_, err := tb.instance.Send(msg)
	if err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}
	return nil
}
