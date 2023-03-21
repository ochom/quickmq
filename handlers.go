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

func publish(x *channel) gin.HandlerFunc {
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
			log.Println("Error unmarshaling: ", err.Error())
			return
		}

		if req.Delay <= CronJobInterval {
			item := models.NewItem(req.Queue, req.Data, req.Delay)
			x.publish(item)
			return
		}

		// if delay is greater than cron time, add queue and items to database
		item := models.NewItem(req.Queue, req.Data, req.Delay)
		if err := x.repo.SaveItem(item); err != nil {
			log.Println("Error adding to repo: ", err.Error())
		}
	}
}

func consume(channel *channel) gin.HandlerFunc {
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

		messages := channel.consume(queueName)
		for msg := range messages {
			if err := ws.WriteMessage(websocket.BinaryMessage, msg); err != nil {
				log.Printf("Error while writing to websocket: %v", err)
				return
			}
		}
	}
}

func getQueues(x *channel) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(200, x.getQueues())
	}
}
