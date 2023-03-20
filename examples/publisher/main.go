package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
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

	wg := &sync.WaitGroup{}
	wg.Add(20)

	for i := 0; i < 20; i++ {
		go func(count int) {
			defer wg.Done()

			b, err := json.Marshal(&examples.Message{Body: text})
			if err != nil {
				panic(err)
			}

			publish(b, delay)
		}(i)
	}

	wg.Wait()
	log.Println("All messages published")
}
