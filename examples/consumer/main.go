package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"

	"github.com/ochom/quickmq/clients"
	"github.com/ochom/quickmq/examples"
)

func main() {
	workerFunc := func(msg []byte) error {
		var message examples.Message
		if err := json.Unmarshal(msg, &message); err != nil {
			return err
		}

		bN, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return err
		}

		if bN.Int64()%2 != 0 {
			return fmt.Errorf("error processing message")
		}

		log.Printf("Got message: %s ", message.Body)
		return nil
	}

	quit := make(chan bool)
	workers := 5
	for i := 0; i < workers; i++ {
		go func(id int) {
			consumer := clients.NewConsumer("ws://localhost:3456", "test-queue", workerFunc, workers)
			if err := consumer.Consume(quit); err != nil {
				log.Println("Error consuming: ", err.Error())
			}
		}(i)
	}

	log.Println("Press CTRL-C to exit")
	<-quit
}
