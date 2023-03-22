package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"

	"github.com/ochom/quickmq/clients"
	"github.com/ochom/quickmq/examples"
)

func main() {
	i := 0
	workerFunc := func(msg []byte) {
		var message examples.Message
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Println("Error unmarshaling message: ", err.Error())
			return
		}

		log.Printf("Got message [%d]: %s ", i, message.Body)
		i++
	}

	go func() {
		consumer := clients.NewConsumer("ws://localhost:3456", "test-queue")
		if err := consumer.Consume(workerFunc); err != nil {
			log.Println("Error consuming: ", err.Error())
		}
	}()

	log.Println("Waiting for messages. To exit press CTRL+C")

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)

	// Wait for a message on the exit quickMQ
	<-exit
}
