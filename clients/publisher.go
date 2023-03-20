package clients

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ochom/quickmq/dto"

	"github.com/ochom/gttp"
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

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	url := fmt.Sprintf("%s/publish", p.url)
	res, status, err := gttp.NewRequest(url, headers, data).Post()
	if err != nil {
		return err
	}

	if status != 200 {
		return fmt.Errorf("got status code %d: %s", status, string(res))

	}

	return nil
}
