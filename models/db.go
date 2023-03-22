package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Repo is a struct that holds the database connection
type Repo struct {
	DB *gorm.DB
}

// NewRepo creates a new Repo
func NewRepo() (*Repo, error) {
	db, err := gorm.Open(sqlite.Open("/var/pubsub/data/db.sqlite"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&QueueItem{}); err != nil {
		return nil, err
	}

	return &Repo{DB: db}, nil
}

// SaveItems adds  QueueItems to the database
func (r *Repo) SaveItems(items []*QueueItem) error {
	return r.DB.Create(&items).Error
}

// SaveItem adds a QueueItem to the database
func (r *Repo) SaveItem(item *QueueItem) error {
	return r.DB.Create(&item).Error
}

// GetQueueItems gets all the queues from the database
func (r *Repo) GetQueueItems(delay int64) ([]*QueueItem, error) {

	ids := make([]string, 0)
	var items []*QueueItem
	err := r.DB.Where("delay <= ?", delay).Find(&items).Error
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		ids = append(ids, item.ID)
	}

	if len(ids) > 0 {
		if err := r.DeleteQueueItems(ids); err != nil {
			return nil, err
		}
	}

	return items, nil
}

// GetAll gets all the queues from the database
func (r *Repo) GetAll() ([]*QueueItem, error) {
	var items []*QueueItem

	if err := r.DB.Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

// DeleteQueueItems deletes all the queues from the database
func (r *Repo) DeleteQueueItems(ids []string) error {
	err := r.DB.Exec("DELETE FROM queue_items WHERE id IN ?", ids).Error
	if err != nil {
		return err
	}

	return nil
}
