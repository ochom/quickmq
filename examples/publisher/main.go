package main

import (
	"encoding/json"
	"log"
	"ochom/pubsub/clients"
	"ochom/pubsub/examples"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		panic("no arguments")
	}

	publish := func(b []byte, delay int) {
		pubslisher := clients.NewPublisher("http://localhost:8081", "test-queue")
		if err := pubslisher.PublishWithDelay(b, time.Duration(delay)*time.Second); err != nil {
			panic(err)
		}
	}

	// delay is the last argument
	delay, err := strconv.Atoi(args[len(args)-1])
	if err != nil {
		panic(err)
	}

	message := examples.Message{
		Body: strings.Join(args[0:len(args)-1], " "),
	}

	b, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 100; i++ {
		publish(b, delay+i)
	}

	log.Println("All messages published")
}
