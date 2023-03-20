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
}

// NewConsumer creates a new Consumer
func NewConsumer(url string, queueName string, workerFunc WorkerFunc) *Consumer {
	return &Consumer{
		host:       url,
		queueName:  queueName,
		workerFunc: workerFunc,
	}
}

// Consume consumes a message from a queue
func (c *Consumer) Consume() error {
	msgs := make(chan []byte)
	go c.getMessages(msgs)

	for msg := range msgs {
		if err := c.workerFunc(msg); err != nil {
			log.Println("Error processing message: ", err)
			c.reQueue(msg)
		}
	}

	return nil
}

// getMessages consumes a message from a queue
func (c *Consumer) getMessages(deliveries chan []byte) {
	cl, _, err := websocket.DefaultDialer.Dial(c.host+"/consume?queue="+c.queueName, nil)
	if err != nil {
		log.Println("Error dialing: ", err.Error())
		return
	}

	defer cl.Close()

	for {
		_, data, err := cl.ReadMessage()
		if err != nil {
			log.Println("Error reading message: ", err.Error())
			return
		}

		deliveries <- data
	}
}

// reQueue requeues a message
func (c *Consumer) reQueue(data []byte) {
	cl, _, err := websocket.DefaultDialer.Dial(c.host+"/requeue", nil)
	if err != nil {
		log.Println("Error dialing: ", err.Error())
		return
	}

	defer cl.Close()

	if err := cl.WriteMessage(websocket.BinaryMessage, data); err != nil {
		log.Println("Error writing message: ", err.Error())
		return
	}

}
