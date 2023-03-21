package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// create a cron job that runs every 30 minutes
func cron(mq *quickMQ, stop chan os.Signal) {
	ticker := time.NewTicker(CronJobInterval)
	for {
		select {
		case <-ticker.C:
			until := time.Until(time.Now().Add(CronJobInterval))
			items, err := mq.repo.GetQueueItems(until)
			if err != nil {
				log.Println("Error getting ready items: ", err.Error())
				continue
			}

			for _, item := range items {
				mq.publish(item)
			}
		case <-stop:
			ticker.Stop()
			return
		}
	}
}

func main() {
	mq, err := newQuickMQ()
	if err != nil {
		log.Fatalf("Error while creating quickMQ: %v", err)
	}

	server := gin.New()
	server.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{
			"/ping",
			"/publish",
			"/consume",
		},
	}))
	server.Use(gin.Recovery())

	// gin hide paths

	server.GET("/ping", ping())
	server.Any("/publish", publish(mq))
	server.Any("/consume", consume(mq))

	// api
	server.GET("/api/queues", getQueues(mq))

	go func() {
		if err := server.Run(":8080"); err != nil {
			log.Fatalf("Error while running server: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go cron(mq, stop)
	<-stop

	mq.kill()
}
