package model

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID         int        `gorm:"type:int;primary_key" json:"id"`
	Title      string     `gorm:"type:varchar(50);not null" json:"title"`
	Expiration *time.Time `gorm:"type:datetime" json:"expiration"`
	ChatID     int64      `gorm:"type:bigint" json:"chat_id"`
	CreatedAt  time.Time  `gorm:"type:timestamp" json:"created_at"`
	*gorm.Model
}
