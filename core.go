package main

import (
	"sync"

	"github.com/ochom/quickmq/models"
)

// quickMQ  is a  in-memory queue
type quickMQ struct {
	instant map[string]chan []byte
	mutext  sync.Mutex
}

// newQuickMQ creates a new quickMQ
func newQuickMQ() (*quickMQ, error) {
	return &quickMQ{
		instant: make(map[string]chan []byte),
	}, nil
}

// open opens a queue
func (c *quickMQ) open(queue string) (chan []byte, error) {
	c.mutext.Lock()
	defer c.mutext.Unlock()

	q, ok := c.instant[queue]
	if !ok {
		q = make(chan []byte, 1)
		c.instant[queue] = q
	}

	return q, nil
}

// Add adds a QueueItem to the queue
func (c *quickMQ) publish(item *models.QueueItem) {
	q, err := c.open(item.QueueName)
	if err != nil {
		return
	}

	q <- item.Data
}

// consume returns a channel that will return the next item in the queue
func (c *quickMQ) consume(queue string) chan []byte {
	q, err := c.open(queue)
	if err != nil {
		return nil
	}

	return q
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
