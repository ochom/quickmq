package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"sync/atomic"

	"github.com/ochom/quickmq/clients"
	"github.com/ochom/quickmq/examples"
)

func main() {
	i := int32(0)
	workerFunc := func(msg []byte) {
		var message examples.Message
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Println("Error unmarshaling message: ", err.Error())
			return
		}

		log.Printf("Got message [%d]: %s ", i, message.Body)
		atomic.AddInt32(&i, 1)
	}

	for i := 0; i < 10; i++ {
		go func(worker int) {
			c := clients.NewConsumer("ws://localhost:3456", "test-queue")
			if err := c.Consume(workerFunc); err != nil {
				log.Printf("[%s] workerID [%d] Error consuming: %v\n", c.GetQueueName(), worker, err)
			}
		}(i)
	}

	log.Println("Waiting for messages. To exit press CTRL+C")

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)

	// Wait for a message on the exit quickMQ
	<-exit
}
