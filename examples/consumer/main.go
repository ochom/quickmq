package main

import (
	"log"
	"ochom/pubsub/clients"
)

func main() {

	workChannel := make(chan []byte)
	go func() {
		for msg := range workChannel {
			log.Println("Got message: ", string(msg))
		}
	}()

	consumer := clients.NewConsumer("http://localhost:8080", "test-queue")
	if err := consumer.Consume(workChannel); err != nil {
		panic(err)
	}

	log.Println("Press CTRL-C to exit")
	quit := make(chan bool)
	<-quit
}
