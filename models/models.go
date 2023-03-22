package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// Queue is a struct that holds the name of the queue and the items in the queue
type Queue struct {
	Name  string       `json:"name,omitempty"`
	Items []*QueueItem `json:"items,omitempty"`
}

// NewQueue creates a new Queue
func NewQueue(name string) *Queue {
	return &Queue{
		Name: name,
	}
}

// QueueItem is a struct that holds the data for a queue item
type QueueItem struct {
	ID        string         `gorm:"primaryKey"`
	QueueName string         `json:"queue_name"`
	Data      datatypes.JSON `json:"data"`
	Delay     int64          `json:"delay"`
}

// NewItem creates a new QueueItem
func NewItem(qName string, data []byte, delay time.Duration) *QueueItem {
	return &QueueItem{
		ID:        uuid.New().String(),
		QueueName: qName,
		Data:      data,
		Delay:     time.Now().Add(delay).Unix(),
	}
}
