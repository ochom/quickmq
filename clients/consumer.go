package clients

import (
	"log"

	"github.com/gorilla/websocket"
)

// WorkerFunc is a function that processes a message
type WorkerFunc func([]byte) error

// Consumer is a struct that holds the name of the queue and the items in the queue
type Consumer struct {
	host       string
	queueName  string
	workerFunc WorkerFunc
	workers    int
}

// NewConsumer creates a new Consumer
func NewConsumer(url string, queueName string, workerFunc WorkerFunc, workers int) *Consumer {
	return &Consumer{
		host:       url,
		queueName:  queueName,
		workerFunc: workerFunc,
		workers:    workers,
	}
}

// GetQueueName returns the name of the queue
func (c *Consumer) GetQueueName() string {
	return c.queueName
}

// GetWorkers returns the number of workers
func (c *Consumer) GetWorkers() int {
	return c.workers
}

// Consume consumes a message from a queue
func (c *Consumer) Consume() error {
	msgs := c.getMessages()

	for msg := range msgs {
		if err := c.workerFunc(msg); err != nil {
			c.reQueue(msg)
		}
	}

	return nil
}

// getMessages consumes a message from a queue
func (c *Consumer) getMessages() <-chan []byte {
	cl, _, err := websocket.DefaultDialer.Dial(c.host+"/consume?queue="+c.queueName, nil)
	if err != nil {
		log.Println("Error dialing: ", err.Error())
		return nil
	}

	deliveries := make(chan []byte)
	go func() {

		for {
			_, data, err := cl.ReadMessage()
			if err != nil {
				log.Println("Error reading message: ", err.Error())
				break
			}

			deliveries <- data
		}

		close(deliveries)
	}()

	return deliveries
}

// reQueue re-queues a message
func (c *Consumer) reQueue(data []byte) {
	p := NewPublisher(c.host, c.queueName)
	if err := p.Publish(data); err != nil {
		log.Println("Error requeuing message: ", err.Error())
		return
	}
}
