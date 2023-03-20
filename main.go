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

func main() {
	stop := make(chan os.Signal, 1)
	channel := NewChannel(stop)

	server := gin.New()
	server.Use(gin.Logger())
	server.Use(gin.Recovery())

	server.GET("/ping", ping())
	server.Any("/publish", publish(channel))
	server.Any("/consume", consume(channel))

	go func() {
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
