package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"ochom/pubsub/dto"
	"ochom/pubsub/models"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	stop := make(chan os.Signal, 1)
	channel := NewChannel(stop)

	server := gin.New()
	server.Use(gin.Logger())
	server.Use(gin.Recovery())

	server.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	server.POST("/publish", func(c *gin.Context) {
		queue := c.Query("queue")
		if queue == "" {
			c.JSON(400, gin.H{"error": "queue is required"})
			return
		}

		data, err := c.GetRawData()
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		req := &dto.PublishRequest{}
		if err := json.Unmarshal(data, req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		q := models.NewQueue(queue)
		item := models.NewItem(q.ID, req.Data, req.Delay)

		channel.Add(q, item)

		c.JSON(200, gin.H{"message": "ok"})
	})

	server.GET("/consume", func(c *gin.Context) {
		queue := c.Query("queue")
		if queue == "" {
			c.JSON(400, gin.H{"error": "queue is required"})
			return
		}

		stream := make(chan []byte)
		go channel.Consume(models.QueueName(queue), stream)
		c.Stream(func(w io.Writer) bool {
			for {
				select {
				case <-stop:
					return false
				case item := <-stream:
					c.JSON(200, item)
					c.Writer.Flush()
				}
			}
		})
	})

	go func() {
		log.Println("Starting server...")
		if err := server.Run(":8080"); err != nil {
			log.Fatalf("Error while running server: %v", err)
		}
	}()

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
