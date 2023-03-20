package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/ochom/quickmq/models"
)

// Queue is a struct that holds the name of the queue and the items in the queue
type Channel struct {
	repo   *models.Repo
	queues map[string][]*models.QueueItem
	stop   chan os.Signal
	mutext sync.Mutex
}

// NewChannel creates a new Channel
func NewChannel(stop chan os.Signal) *Channel {
	repo, err := models.NewRepo()
	if err != nil {
		panic(err)
	}

	return &Channel{
		queues: make(map[string][]*models.QueueItem),
		repo:   repo,
		stop:   stop,
	}
}

// Add adds a QueueItem to the queue
func (c *Channel) Add(queue *models.Queue, item *models.QueueItem) {
	c.mutext.Lock()
	defer c.mutext.Unlock()

	q, ok := c.queues[queue.Name]
	if !ok {
		q = make([]*models.QueueItem, 0)
	}

	q = append(q, item)
	c.queues[queue.Name] = q
}

// Get gets the next QueueItem from the queue that is ready to be sent
func (c *Channel) Get(qName string) *models.QueueItem {
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

// Consume consumes the queue and sends the items to the given channel
// if stop signal is received, the function returns
func (c *Channel) Consume(queue string, ch chan<- []byte) {
	for {
		select {
		case <-c.stop:
			return
		default:
			item := c.Get(queue)
			if item != nil {
				ch <- item.Data
			}
		}
	}
}

// Stop stops the channel and writes all items to disk
func (c *Channel) Stop(ctx context.Context) error {

	updateData := func() {
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

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context timed out")
		default:
			updateData()
			return nil
		}
	}
}

// GetQueues returns a list of queues
func (c *Channel) GetQueues() []*models.Queue {
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
