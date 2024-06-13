package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/ochom/gutils/logs"
	"github.com/ochom/quickmq/src/api"
	"github.com/ochom/quickmq/src/app"
)

func main() {
	coreServer := app.New()
	webServer := api.New()

	// run core
	go func() {
		if err := coreServer.Listen(":6321"); err != nil {
			panic(err)
		}
	}()

	// run api and web
	go func() {
		if err := webServer.Listen(":16321"); err != nil {
			panic(err)
		}
	}()

	// go run consumer daemon
	go func() {
		stopSignal := make(chan bool, 1)
		logs.Info("starting consumers daemon")
		app.Consume(stopSignal)

		stopSignal <- true
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// shutdown core server
	if err := coreServer.ShutdownWithContext(ctx); err != nil {
		panic(err)
	}

	// shutdown api server
	if err := webServer.ShutdownWithContext(ctx); err != nil {
		panic(err)
	}

	logs.Info("Server shutdown correctly")
}
