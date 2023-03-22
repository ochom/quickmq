package main

import (
	"encoding/json"
	"log"

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
			log.Println("Published to queue")
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

		defer ws.Close()

		messages := mq.consume(queueName)
		for msg := range messages {
			if err := ws.WriteMessage(websocket.BinaryMessage, msg); err != nil {
				log.Printf("Error while writing to websocket: %v", err)
				return
			}
		}
	}
}

func getQueues(mq *quickMQ, repo *models.Repo) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		res, err := mq.getQueues(repo)
		if err != nil {
			ctx.String(200, "No queues found")
			return
		}

		ctx.JSON(200, res)
	}
}
