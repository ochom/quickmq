package main

import (
	"encoding/json"
	"log"

	"github.com/ochom/quickmq/clients"
	"github.com/ochom/quickmq/examples"
)

func main() {

	workerFunc := func(msg []byte) error {
		var message examples.Message
		if err := json.Unmarshal(msg, &message); err != nil {
			return err
		}

		log.Printf("Got message: %s ", message.Body)
		return nil
	}

	workers := 5
	for i := 0; i < workers; i++ {
		go func(id int) {
			consumer := clients.NewConsumer("ws://localhost:8080", "test-queue", workerFunc)
			if err := consumer.Consume(); err != nil {
				log.Println("Error consuming: ", err.Error())
			}
		}(i)
	}

	log.Println("Press CTRL-C to exit")
	quit := make(chan bool)
	<-quit
}
