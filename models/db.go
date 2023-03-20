package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Repo is a struct that holds the database connection
type Repo struct {
	DB *gorm.DB
}

// NewRepo creates a new Repo
func NewRepo() (*Repo, error) {
	db, err := gorm.Open(sqlite.Open("pubsub.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&Queue{}, &QueueItem{}); err != nil {
		return nil, err
	}

	return &Repo{DB: db}, nil
}

// AddQueue adds a queue to the database
func (r *Repo) AddQueue(queue *Queue) error {
	return r.DB.Create(&queue).Error
}

// AddQueueItem adds a queue item to the database
func (r *Repo) AddQueueItems(items []*QueueItem) error {
	return r.DB.Create(&items).Error
}

// GetQueues gets all the queues from the database
func (r *Repo) GetQueues() ([]Queue, error) {
	var queues []Queue
	err := r.DB.Find(&queues).Error
	if err != nil {
		return nil, err
	}

	for i := range queues {
		queues[i].Items, err = r.GetQueueItems(queues[i].ID)
		if err != nil {
			return nil, err
		}
	}

	return queues, nil
}

// GetQueueItems gets all the items from a queue
func (r *Repo) GetQueueItems(queueID string) ([]QueueItem, error) {
	var items []QueueItem
	err := r.DB.Where("queue_id = ?", queueID).Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

// DeleteQueues deletes all the queues from the database
func (r *Repo) DeleteQueues() error {
	err := r.DB.Exec("DELETE FROM queues").Error
	if err != nil {
		return err
	}

	err = r.DB.Exec("DELETE FROM queue_items").Error
	if err != nil {
		return err
	}

	return nil
}
