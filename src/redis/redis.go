package redis

import (
	"context"
	"fmt"
	"log"
	"os"
	"remind_me/src/db"
	"remind_me/src/utils"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	instance *redis.Client
}

func StartRedis(dbInstance *db.DB, SendTelegramMessage func(chatID int64, message string) error) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     utils.GetEnv("REDIS_HOST"),
		Password: utils.GetEnv("REDIS_PASSWORD"),
		DB:       0,
	})
	_, err := client.Do(context.Background(), "CONFIG", "SET", "notify-keyspace-events", "KEA").Result()
	if err != nil {
		fmt.Printf("unable to set keyspace events %v", err.Error())
		os.Exit(1)
	}

	pubsub := client.PSubscribe(context.Background(), "__keyevent@0__:expired")
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
				expiredKeyInt, err := strconv.Atoi(expiredKey)
				if err != nil {
					panic(err)
				}

				result, err := dbInstance.GetReminderById(expiredKeyInt)
				if err != nil {
					panic("error trying to get reminder from db: " + err.Error())
				}

				message := fmt.Sprintf("Reminder:  '%s' has expired", result.Title)
				utils.SuccessLog("Message sent!")
				err = SendTelegramMessage(result.ChatID, message)
				if err != nil {
					panic("error trying to send message: " + err.Error())
				}
			}
		}
	}()

	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Error connecting to Redis:", err)
	}
	utils.SuccessLog("Successful connection to Redis")

	return &Redis{instance: client}, nil
}

func (r *Redis) Set(key string, value string, expiration time.Duration) error {
	status := r.instance.Set(context.Background(), key, value, expiration)
	if status.Err() != nil {
		return fmt.Errorf("error setting value in Redis: %w", status.Err())
	}
	return nil
}
