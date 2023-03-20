package clients

import (
	"fmt"
	"net/http"
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

// Consume consumes a message from a queue
func (c *Consumer) Consume(workerChan chan []byte) error {
	headers := map[string]string{
		"Accept":      "application/json",
		"requestType": "stream",
	}

	url := fmt.Sprintf("%s/consume?queue=%s", c.host, c.queueName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("got status code %d: %s", res.StatusCode, res.Status)
	}

	for {
		data := make([]byte, 1024)
		n, err := res.Body.Read(data)
		if err != nil {
			return err
		}

		workerChan <- data[:n]
	}
}
