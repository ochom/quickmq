package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// QueueName is a string that holds the name of the queue
type QueueName string

// Queue is a struct that holds the name of the queue and the items in the queue
type Queue struct {
	ID    string      `gorm:"primaryKey"`
	Name  QueueName   `gorm:"uniqueIndex"`
	Items []QueueItem `gorm:"-"` // ignore this field
}

// NewQueue creates a new Queue
func NewQueue(name string) *Queue {
	return &Queue{
		ID:   uuid.New().String(),
		Name: QueueName(name),
	}
}

// AfterFind is a hook that is called after a record is found
func (q *Queue) AfterFind(tx *gorm.DB) (err error) {
	var items []QueueItem
	tx.Where("queue_id = ?", q.ID).Find(&items)
	q.Items = items
	return
}

type QueueItem struct {
	ID         string `gorm:"primaryKey"`
	QueueID    string `gorm:"index"`
	StringData string `json:"data"`
	SendAt     int64  `json:"send_at"`
}

// NewItem creates a new QueueItem
func NewItem(queueID string, data []byte, delay time.Duration) *QueueItem {
	return &QueueItem{
		ID:         uuid.New().String(),
		QueueID:    queueID,
		StringData: string(data),
		SendAt:     time.Now().Add(delay).Unix(),
	}
}

// GetData returns the data as a byte array
func (i *QueueItem) GetData() []byte {
	return []byte(i.StringData)
}
