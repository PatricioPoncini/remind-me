package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	instance *redis.Client
}

func StartRedis(dbInstance *DB, bot *tgbotapi.BotAPI) (*Redis, error) {
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
					fmt.Println("Error:", err)
					return
				}

				resultJSON, err := json.MarshalIndent(result, "", "  ")
				if err != nil {
					log.Fatalf("Error marshaling result to JSON: %v", err)
				}

				// logger
				fmt.Println(string(resultJSON))
				message := fmt.Sprintf("Reminder:  '%s' has expired", result.Title)
				err = SendTelegramMessage(bot, result.ChatID, message)
				if err != nil {
					fmt.Println("Error sending message:", err)
					return
				}
			}
		}
	}()

	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Error connecting to Redis:", err)
	}
	fmt.Println("\033[32m- Successful connection to Redis\033[0m")

	return &Redis{instance: client}, nil
}

func (r *Redis) Set(key string, value string, expiration time.Duration) error {
	status := r.instance.Set(context.Background(), key, value, expiration)
	if status.Err() != nil {
		return fmt.Errorf("error setting value in Redis: %w", status.Err())
	}
	return nil
}
