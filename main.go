package main

import (
	"remind_me/src/bot"
	"remind_me/src/controller"
	"remind_me/src/db"
	"remind_me/src/model"
	"remind_me/src/repo"
	"remind_me/src/use_case"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func main() {
	mysqlDB := db.ConnectionMySQLDB("root:password@tcp(localhost:3306)/remind_me_db")
	mysqlDB.AutoMigrate(&model.Message{})
	redisDB := db.ConnectionRedisDB()

	messageController := startServer(mysqlDB, redisDB)
	telegramBot := bot.StartBot(messageController)

	db.ListenForExpiredKeys(redisDB, messageController, telegramBot)
}

func startServer(db *gorm.DB, rdb *redis.Client) *controller.MessageController {
	messageRepo := repo.NewMessageRepo(db, rdb)
	messageUseCase := use_case.NewMessageUseCase(messageRepo)
	messageController := controller.NewMessageControler(messageUseCase)
	return messageController
}
