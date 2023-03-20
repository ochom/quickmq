package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/ochom/quickmq/models"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// create a cron job that runs every 30 minutes
func cron(channel *Channel) {
	ticker := time.NewTicker(CronJobInterval)
	for {
		select {
		case <-ticker.C:
			until := time.Until(time.Now().Add(CronJobInterval))
			items, err := channel.repo.GetQueueItems(until)
			if err != nil {
				log.Println("Error getting ready items: ", err.Error())
				continue
			}

			for _, item := range items {
				q := models.NewQueue(item.QueueName)
				channel.Add(q, item)
			}
		case <-channel.stop:
			ticker.Stop()
			return
		}
	}
}

func main() {
	stop := make(chan os.Signal, 1)
	channel := NewChannel(stop)

	server := gin.New()
	server.Use(gin.Logger())
	server.Use(gin.Recovery())

	server.GET("/ping", ping())
	server.Any("/publish", publish(channel))
	server.Any("/consume", consume(channel))

	// api
	server.GET("/api/queues", getQueues(channel))

	go func() {
		if err := server.Run(":8080"); err != nil {
			log.Fatalf("Error while running server: %v", err)
		}
	}()

	go cron(channel)

	signal.Notify(stop, os.Interrupt)
	<-stop
	// wait for 5 minutes before exiting
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	log.Println("Exiting...")
	if err := channel.Stop(ctx); err != nil {
		log.Fatalf("Error while exiting: %v", err)
	}

	log.Println("Exited successfully")
}
