package main

import (
	"context"
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
func cron(channel *channel) {
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
				channel.publish(item)
			}
		case <-channel.stop:
			ticker.Stop()
			return
		}
	}
}

func main() {
	stop := make(chan os.Signal, 1)
	x := newChannel(stop)

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
	server.Any("/publish", publish(x))
	server.Any("/consume", consume(x))

	// api
	server.GET("/api/queues", getQueues(x))

	go func() {
		if err := server.Run(":8080"); err != nil {
			log.Fatalf("Error while running server: %v", err)
		}
	}()

	go cron(x)

	signal.Notify(stop, os.Interrupt)
	<-stop
	// wait for 5 minutes before exiting
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	log.Println("Exiting...")
	if err := x.kill(ctx); err != nil {
		log.Fatalf("Error while exiting: %v", err)
	}

	log.Println("Exited successfully")
}
