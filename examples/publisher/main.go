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
		p := clients.NewPublisher("ws://localhost:3456", "test-queue")
		if err := p.PublishWithDelay(b, time.Duration(delay)*time.Second); err != nil {
			panic(err)
		}
	}

	// delay is the last argument
	delay, err := strconv.Atoi(args[len(args)-1])
	if err != nil {
		panic(err)
	}

	workSize := int32(100000)

	jobs := make(chan []byte, workSize)
	results := make(chan bool, workSize)

	for w := 1; w <= 100; w++ {
		go func(worker int) {
			for j := range jobs {
				publish(j, delay)
				results <- true
			}
		}(w)
	}

	var i int32
	for i = 0; i < workSize; i++ {
		message := examples.Message{
			Body: text + " " + strconv.Itoa(int(i)),
		}

		b, err := json.Marshal(message)
		if err != nil {
			panic(err)
		}

		jobs <- b
	}

	close(jobs)

	for a := int32(0); a < workSize; a++ {
		<-results
	}
	log.Println("All messages published")
}
