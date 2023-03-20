package clients

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/ochom/quickmq/dto"
)

// Publisher is a struct that holds the name of the queue and the items in the queue
type Publisher struct {
	url       string
	queueName string
}

// NewPublisher creates a new Publisher
func NewPublisher(url string, queueName string) *Publisher {
	return &Publisher{
		url:       url,
		queueName: queueName,
	}
}

// PublishWithDelay publishes a message to a queue with a delay
func (p *Publisher) PublishWithDelay(d []byte, delay time.Duration) error {
	return p.publish(d, delay)
}

// Publish publishes a message to a queue
func (p *Publisher) Publish(d []byte) error {
	return p.publish(d, 0)
}

// Publish publishes a message to a queue
func (p *Publisher) publish(d []byte, delay time.Duration) error {
	payload := &dto.PublishRequest{
		Queue: p.queueName,
		Data:  d,
		Delay: delay,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/publish", p.url)

	cl, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Println("Error dialing: ", err.Error())
		return err
	}

	defer cl.Close()

	if err := cl.WriteMessage(websocket.BinaryMessage, data); err != nil {
		log.Println("Error writing message: ", err.Error())
		return err
	}

	return nil
}
