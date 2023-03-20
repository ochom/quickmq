package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// Queue is a struct that holds the name of the queue and the items in the queue
type Queue struct {
	Name  string       `gorm:"uniqueIndex"`
	Items []*QueueItem `gorm:"-"` // ignore this field
}

// NewQueue creates a new Queue
func NewQueue(name string) *Queue {
	return &Queue{
		Name: name,
	}
}

type QueueItem struct {
	ID        string         `gorm:"primaryKey"`
	QueueName string         `json:"queue_name"`
	Data      datatypes.JSON `json:"data"`
	SendAt    int64          `json:"send_at"`
}

// NewItem creates a new QueueItem
func NewItem(qName string, data []byte, delay time.Duration) *QueueItem {
	return &QueueItem{
		ID:        uuid.New().String(),
		QueueName: qName,
		Data:      data,
		SendAt:    time.Now().Add(delay).Unix(),
	}
}
