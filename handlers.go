package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/ochom/quickmq/dto"
	"github.com/ochom/quickmq/models"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func ping() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	}
}

func publish(channel *Channel) gin.HandlerFunc {
	return func(c *gin.Context) {
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		defer ws.Close()

		_, data, err := ws.ReadMessage()
		if err != nil {
			log.Println("Error reading message: ", err.Error())
			return
		}

		req := &dto.PublishRequest{}
		if err := json.Unmarshal(data, req); err != nil {
			log.Println("Error unmarshaling: ", err.Error())
			return
		}

		q := models.NewQueue(req.Queue)
		item := models.NewItem(q.ID, req.Data, req.Delay)

		channel.Add(q, item)
	}
}

func consume(channel *Channel) gin.HandlerFunc {
	return func(c *gin.Context) {
		queue := c.Query("queue")
		if queue == "" {
			c.JSON(400, gin.H{"error": "queue is required"})
			return
		}

		queueName := models.QueueName(queue)

		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		defer ws.Close()

		for {
			select {
			case <-channel.stop:
				return
			default:
				item := channel.Get(queueName)
				if item == nil {
					time.Sleep(1 * time.Second)
					continue
				}

				if err := ws.WriteMessage(websocket.BinaryMessage, item.Data); err != nil {
					log.Printf("Error while writing to websocket: %v", err)
					return
				}
			}
		}
	}
}

func getQueues(channel *Channel) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, channel.GetQueues())
	}
}

func getMessages(channel *Channel) gin.HandlerFunc {
	return func(c *gin.Context) {
		queue := c.Param("queueName")
		if queue == "" {
			c.JSON(400, gin.H{"error": "queue is required"})
			return
		}
		c.JSON(200, channel.GetQueueItems(models.QueueName(queue)))
	}
}
