package bot

import (
	"fmt"
	"log"
	"remind_me/src/controller"
	"remind_me/src/model"
	"remind_me/src/utils"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func StartBot(msgController *controller.MessageController) *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(utils.GetEnv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}

	u := tgbotapi.NewUpdate(0)
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatalf("Error getting updates: %v", err)
	}

	for update := range updates {
		isValidCommand := strings.HasPrefix(update.Message.Text, utils.RemindCommand)
		if !isValidCommand {
			fmt.Println("Error. Invalid message")
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "The only command to use for now is '/r' to set a reminder. The message you sent is not valid.")
			_, err := bot.Send(msg)
			if err != nil {
				panic("Error trying to send message: " + err.Error())
			}
		} else {
			title, notifyTime, duration, err := utils.ParseMessage(update.Message.Text) // remove duration
			if err != nil {
				panic("Error trying to parse message: " + err.Error())
			}

			if update.Message == nil { // Ignore any non-Message updates.
				continue
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Reminder setted: '"+title+"'")

			message := model.Message{
				Title:      title,
				Expiration: &notifyTime,
				ChatID:     update.Message.Chat.ID,
				CreatedAt:  time.Now(), // Establecer la fecha actual
			}
			err = msgController.CreateMessage(message, duration)

			if err != nil {
				panic("Error trying to insert new reminder in db: " + err.Error())
			}

			_, err = bot.Send(msg)
			if err != nil {
				panic("Error trying to send message: " + err.Error())
			}
		}
	}

	return bot
}

func SendTelegramMessage(bot *tgbotapi.BotAPI, chatID int64, message string) error {
	msg := tgbotapi.NewMessage(chatID, message)
	_, err := bot.Send(msg)
	if err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}
	return nil
}
