package clients

import (
	"log"

	"github.com/gorilla/websocket"
)

// Consumer is a struct that holds the name of the queue and the items in the queue
type Consumer struct {
	host      string
	queueName string
}

// NewConsumer creates a new Consumer
func NewConsumer(url string, queueName string) *Consumer {
	return &Consumer{
		host:      url,
		queueName: queueName,
	}
}

// GetQueueName returns the name of the queue
func (c *Consumer) GetQueueName() string {
	return c.queueName
}

// Consume consumes a message from a queue
func (c *Consumer) Consume(workerFunc func([]byte)) error {
	cl, res, err := websocket.DefaultDialer.Dial(c.host+"/consume?queue="+c.queueName, nil)
	if err != nil {
		log.Println("Error dialing: ", err.Error())
		return nil
	}

	if res.StatusCode != 101 {
		log.Println("Error dialing: code:", res.StatusCode, res.Status)
		return nil
	}

	defer cl.Close()

	for {
		_, data, err := cl.ReadMessage()
		if err != nil {
			log.Println("Error reading message: ", err.Error())
			break
		}

		workerFunc(data)
	}

	return nil
}
