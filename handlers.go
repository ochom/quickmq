package main

import (
	"github.com/ochom/quickmq/models"

	"github.com/gin-gonic/gin"
)

func ping() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
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
