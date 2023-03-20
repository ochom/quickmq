package main

import (
	"encoding/json"
	"log"
	"ochom/pubsub/clients"
	"ochom/pubsub/examples"
)

func main() {

	consumer := clients.NewConsumer("http://localhost:8081", "test-queue")
	workChannel := make(chan []byte)
	go func() {
		for msg := range workChannel {
			var message examples.Message
			if err := json.Unmarshal(msg, &message); err != nil {
				log.Println("Error unmarshalling message: ", err)
			} else {
				log.Println("Got message: ", message.Body)
			}
		}
	}()

	go func() {
		if err := consumer.Consume(workChannel); err != nil {
			panic(err)
		}
	}()

	log.Println("Press CTRL-C to exit")
	quit := make(chan bool)
	<-quit
}
