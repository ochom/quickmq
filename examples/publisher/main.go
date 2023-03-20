package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ochom/quickmq/clients"
	"github.com/ochom/quickmq/examples"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		panic("no arguments")
	}

	text := strings.Join(args[0:len(args)-1], " ")

	publish := func(b []byte, delay int) {
		pubslisher := clients.NewPublisher("http://localhost:8080", "test-queue")
		if err := pubslisher.PublishWithDelay(b, time.Duration(delay)*time.Second); err != nil {
			panic(err)
		}
	}

	// delay is the last argument
	delay, err := strconv.Atoi(args[len(args)-1])
	if err != nil {
		panic(err)
	}

	for i := 0; i < 100; i++ {
		go func(count int) {
			message := examples.Message{
				Body: text + " " + strconv.Itoa(count),
			}

			b, err := json.Marshal(message)
			if err != nil {
				panic(err)
			}

			publish(b, delay)
		}(i)
	}

	log.Println("All messages published")
	time.Sleep(5 * time.Second)
}
