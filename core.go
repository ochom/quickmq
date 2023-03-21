package main

import (
	"log"
	"sync"
	"time"

	"github.com/ochom/quickmq/models"
)

// quickMQ  is a  in-memory queue
type quickMQ struct {
	repo   *models.Repo
	queues map[string][]*models.QueueItem
	mutext sync.Mutex
}

// newQuickMQ creates a new quickMQ
func newQuickMQ() (*quickMQ, error) {
	repo, err := models.NewRepo()
	if err != nil {
		return nil, err
	}

	return &quickMQ{
		queues: make(map[string][]*models.QueueItem),
		repo:   repo,
	}, nil
}

// Add adds a QueueItem to the queue
func (c *quickMQ) publish(item *models.QueueItem) {
	c.mutext.Lock()
	defer c.mutext.Unlock()

	q, ok := c.queues[item.QueueName]
	if !ok {
		q = make([]*models.QueueItem, 0)
	}

	q = append(q, item)
	c.queues[item.QueueName] = q
}

// getItem returns the next item in the queue
func (c *quickMQ) getItem(qName string) *models.QueueItem {
	c.mutext.Lock()
	defer c.mutext.Unlock()

	q, ok := c.queues[qName]
	if !ok {
		return nil
	}

	if len(q) == 0 {
		return nil
	}

	for i, item := range q {
		if time.Until(time.Unix(item.SendAt, 0)) <= 0 {
			c.queues[qName] = append(q[:i], q[i+1:]...)
			return item
		}
	}

	return nil
}

// consume returns a channel that will return the next item in the queue
func (c *quickMQ) consume(queue string) <-chan []byte {
	ch := make(chan []byte, 1)
	defer close(ch)

	go func() {
		for {
			item := c.getItem(queue)
			if item == nil {
				time.Sleep(1 * time.Second)
				continue
			}

			ch <- item.Data
		}
	}()

	return ch
}

// Stop stops the quickMQ and writes all items to disk
func (c *quickMQ) kill() {
	log.Println("Channel stopped, writing all items to disk")
	for _, items := range c.queues {

		batchSize := 5000
		batches := [][]*models.QueueItem{}
		for i := 0; i < len(items); i += batchSize {
			start := i
			end := i + batchSize
			if end > len(items) {
				end = len(items)
			}

			batches = append(batches, items[start:end])
		}

		for _, batch := range batches {
			if err := c.repo.SaveItems(batch); err != nil {
				panic(err)
			}
		}
	}

	log.Println("Channel stopped, all items written to disk")
}

// getQueues returns all queues and their items
func (c *quickMQ) getQueues() []*models.Queue {
	c.mutext.Lock()
	defer c.mutext.Unlock()

	res := make([]*models.Queue, 0)
	for name := range c.queues {
		q := models.NewQueue(name)
		q.Items = c.queues[name]
		res = append(res, q)
	}

	return res
}
