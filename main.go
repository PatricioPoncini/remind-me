package main

import (
	"log"
	"remind_me/src/bot"
	"remind_me/src/db"
	"remind_me/src/redis"
	"remind_me/src/utils"
)

func main() {
	dbInstance, err := db.NewDB(utils.GetEnv("DB_KEY"))
	if err != nil {
		log.Fatalf("Error creating database instance: %v", err)
		return
	}
	dbInstance.CheckInitialConditions()

	telegramBot, err := bot.NewTelegramBot()
	if err != nil {
		log.Fatalf("Error creating Telegram bot: %v", err)
		return
	}

	redisInstance, err := redis.StartRedis(dbInstance, telegramBot.SendTelegramMessage)
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
		return
	}

	telegramBot.Start(dbInstance, redisInstance)
}
