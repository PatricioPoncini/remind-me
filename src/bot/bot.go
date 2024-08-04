package bot

import (
	"fmt"
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
		utils.ErrorLog("Error getting updates: " + err.Error())
		panic(err)
	}

	for update := range updates {
		utils.InfoLog("New incoming message")
		isValidCommand := strings.HasPrefix(update.Message.Text, utils.RemindCommand)
		title, notifyTime, duration, err := utils.ParseMessage(update.Message.Text)
		if !isValidCommand || err != nil {
			utils.ErrorLog("Error. Invalid command")
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "The only command to use for now is '/r' to set a reminder. The message you sent is not valid.")
			_, err := tb.instance.Send(msg)
			if err != nil {
				panic("Error trying to send message: " + err.Error())
			}
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Reminder set:  '"+title+"'")
			result, err := dbInstance.InsertNewReminder(title, notifyTime, update.Message.Chat.ID)
			if err != nil {
				utils.ErrorLog("Error trying to insert new reminder in db: " + err.Error())
				return
			}

			err = redisInstance.Set(strconv.Itoa(result.ID), "", duration)
			if err != nil {
				utils.ErrorLog("Error setting value in Redis:" + err.Error())
				return
			}

			_, err = tb.instance.Send(msg)
			if err != nil {
				utils.ErrorLog("Error trying to send message:" + err.Error())
				return
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
