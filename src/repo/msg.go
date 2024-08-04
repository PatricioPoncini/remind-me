package repo

import (
	"context"
	"errors"
	"fmt"
	"remind_me/src/domain"
	"remind_me/src/model"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type MessageRepo struct {
	db  *gorm.DB
	rdb *redis.Client
}

func (m *MessageRepo) CreateMessage(message model.Message, duration time.Duration) error {
	if err := m.db.Create(&message).Error; err != nil {
		return errors.New("internal server error: MySQL - cannot create message")
	}

	if err := m.rdb.Set(context.Background(), strconv.Itoa(message.ID), "", duration).Err(); err != nil {
		fmt.Println("Redis error:", err)
		return errors.New("internal server error: Redis - cannot create message")
	}

	return nil
}

func (m *MessageRepo) GetMessageById(id int) (model.Message, error) {
	var message model.Message
	if err := m.db.First(&message, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return message, errors.New("message not found")
		}
		return message, errors.New("internal server error: MySQL - cannot retrieve message")
	}
	return message, nil
}

func NewMessageRepo(db *gorm.DB, rdb *redis.Client) domain.MessageRepo {
	return &MessageRepo{
		db:  db,
		rdb: rdb,
	}
}
