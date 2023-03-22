package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/ochom/quickmq/dto"
	"github.com/ochom/quickmq/models"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// create a cron job that runs every 30 minutes
func cron(mq *quickMQ, repo *models.Repo, stop chan os.Signal) {
	ticker := time.NewTicker(CronJobInterval)
	for {
		select {
		case <-ticker.C:
			items, err := repo.GetQueueItems(time.Now().Add(CronJobInterval).Unix())
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
	repo, err := models.NewRepo()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

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
	server.Any("/publish", publish(mq, repo))
	server.Any("/consume", consume(mq))

	// api
	server.GET("/api/queues", getQueues(mq, repo))

	go func() {
		if err := server.Run(":8080"); err != nil {
			log.Fatalf("Error while running server: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go cron(mq, repo, stop)
	<-stop
}

func publish(mq *quickMQ, repo *models.Repo) gin.HandlerFunc {
	return func(c *gin.Context) {
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		defer func() {
			if err := ws.Close(); err != nil {
				log.Println("Error closing websocket: ", err.Error())
			}
		}()

		_, data, err := ws.ReadMessage()
		if err != nil {
			log.Println("Error reading message: ", err.Error())
			return
		}

		req := &dto.PublishRequest{}
		if err := json.Unmarshal(data, req); err != nil {
			log.Println("Error un-marshaling: ", err.Error())
			return
		}

		if req.Delay == 0 {
			item := models.NewItem(req.Queue, req.Data, req.Delay)
			mq.publish(item)
			return
		}

		// if delay is greater than cron time, add queue and items to database
		item := models.NewItem(req.Queue, req.Data, req.Delay)
		if err := repo.SaveItem(item); err != nil {
			log.Println("Error adding to repo: ", err.Error())
		}
	}
}

func consume(mq *quickMQ) gin.HandlerFunc {
	return func(c *gin.Context) {
		queueName := c.Query("queue")
		if queueName == "" {
			c.JSON(400, gin.H{"error": "queue is required"})
			return
		}

		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		closed := make(chan bool, 1)

		// consume messages from queue and write to websocket but stop when websocket is closed
		go func() {
			for {
				select {
				case msg := <-mq.consume(queueName):
					if err := ws.WriteMessage(websocket.TextMessage, msg); err != nil {
						log.Println("Error writing to websocket: ", err.Error())
						return
					}
				case <-closed:
					return
				}
			}
		}()

		// wait for websocket to close
		for {
			_, _, err := ws.ReadMessage()
			if err != nil {
				closed <- true
				break
			}
		}
	}
}
