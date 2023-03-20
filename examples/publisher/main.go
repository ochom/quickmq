package main

import (
	"ochom/pubsub/clients"
	"os"
	"strings"
	"time"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		panic("no arguments")
	}

	message := strings.Join(args[:len(args)-1], " ")

	// delay is the last argument
	delay, err := time.ParseDuration(args[len(args)-1])
	if err != nil {
		panic(err)
	}

	pubslisher := clients.NewPublisher("http://localhost:8080", "test-queue")
	if err := pubslisher.PublishWithDelay([]byte(message), delay); err != nil {
		panic(err)
	}
}
