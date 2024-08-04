package use_case

import (
	"remind_me/src/domain"
	"remind_me/src/model"
	"time"
)

type MessageUseCase struct {
	messageRepo domain.MessageRepo
}

func (m *MessageUseCase) CreateMessage(message model.Message, duration time.Duration) error {
	err := m.messageRepo.CreateMessage(message, duration)
	return err
}

func (m *MessageUseCase) GetMessageById(id int) (model.Message, error) {
	message, err := m.messageRepo.GetMessageById(id)
	if err != nil {
		panic(err)
	}
	return message, nil
}

func NewMessageUseCase(messageRepo domain.MessageRepo) domain.MessageUseCase {
	return &MessageUseCase{
		messageRepo,
	}
}
