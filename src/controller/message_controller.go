package controller

import (
	"remind_me/src/domain"
	"remind_me/src/model"
	"time"
)

type MessageController struct {
	messageUseCase domain.MessageUseCase
}

func NewMessageControler(messageUseCase domain.MessageUseCase) *MessageController {
	return &MessageController{messageUseCase: messageUseCase}
}

func (m *MessageController) CreateMessage(message model.Message, duration time.Duration) error {
	if err := m.messageUseCase.CreateMessage(message, duration); err != nil {
		panic(err) // handle better this
	}

	// send message
	return nil
}

func (m *MessageController) GetMessageById(id int) (model.Message, error) {
	message, err := m.messageUseCase.GetMessageById(id)
	if err != nil {
		panic(err)
	}
	return message, nil
}
