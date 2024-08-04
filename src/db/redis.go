package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"remind_me/src/controller"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Client   *redis.Client
	Addr     string
	Password string
	DB       int
}

func ConnectionRedisDB() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	_, err := client.Do(context.Background(), "CONFIG", "SET", "notify-keyspace-events", "KEA").Result()
	if err != nil {
		fmt.Printf("unable to set keyspace events %v", err.Error())
		os.Exit(1)
	}

	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Error connecting to Redis:", err)
	}
	fmt.Println("\033[32m- Successful connection to Redis\033[0m")

	return client
}

func ListenForExpiredKeys(rdbClient *redis.Client, messageController *controller.MessageController, telegramBot *tgbotapi.BotAPI) {
	pubsub := rdbClient.PSubscribe(context.Background(), "__keyevent@0__:expired")
	go func() {
		defer pubsub.Close()

		for {
			msg, err := pubsub.ReceiveMessage(context.Background())
			if err != nil {
				if err == redis.ErrClosed {
					fmt.Printf("PubSub connection closed: %s\n", err)
					return
				}
				fmt.Printf("Error receiving message: %s\n", err)
				continue
			}

			if msg != nil && msg.Channel == "__keyevent@0__:expired" {
				expiredKey := msg.Payload
				fmt.Println(expiredKey)
			}
		}
	}()
}
