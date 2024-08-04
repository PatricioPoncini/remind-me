package domain

import (
	"remind_me/src/model"
	"time"
)

type MessageRepo interface {
	CreateMessage(message model.Message, duration time.Duration) error
	GetMessageById(id int) (model.Message, error)
}

type MessageUseCase interface {
	CreateMessage(message model.Message, duration time.Duration) error
	GetMessageById(id int) (model.Message, error)
}
