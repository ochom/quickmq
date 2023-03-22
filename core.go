package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/ochom/quickmq/models"
)

// delivery is the message channel content
type delivery chan []byte

// quickMQ  is a  in-memory queue
type quickMQ struct {
	repo    *models.Repo
	instant map[string]delivery
	mutext  sync.Mutex
}

// newQuickMQ creates a new quickMQ
func newQuickMQ() (*quickMQ, error) {
	return &quickMQ{
		instant: make(map[string]delivery),
	}, nil
}

// Add adds a QueueItem to the queue
func (c *quickMQ) publish(item *models.QueueItem) error {
	c.mutext.Lock()
	defer c.mutext.Unlock()

	if time.Until(time.Unix(item.Delay, 0)) > 0 {
		return fmt.Errorf("item not delayed")
	}

	q, ok := c.instant[item.QueueName]
	if !ok {
		q = make(delivery, 1)
		c.instant[item.QueueName] = q
	}

	q <- item.Data

	return nil
}

// consume returns a channel that will return the next item in the queue
func (c *quickMQ) consume(queue string) delivery {
	c.mutext.Lock()
	defer c.mutext.Unlock()
	return c.instant[queue]
}

// getQueues returns all queues and their items
func (c *quickMQ) getQueues(repo *models.Repo) ([]*models.Queue, error) {
	items, err := repo.GetAll()
	if err != nil {
		return nil, err
	}

	queues := make(map[string][]*models.QueueItem)

	for _, item := range items {
		q, ok := queues[item.QueueName]
		if !ok {
			q = []*models.QueueItem{item}
		}

		q = append(q, item)

		queues[item.QueueName] = q
	}

	res := []*models.Queue{}
	for k, v := range queues {
		q := models.NewQueue(k)
		q.Items = v
		res = append(res, q)

	}

	return res, nil
}
