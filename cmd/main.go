package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/ochom/gutils/logs"
	"github.com/ochom/quickmq/src/app"
)

func main() {
	svr := app.New()
	port := ":16321"
	go func() {
		if err := svr.Listen(port); err != nil {
			panic(err)
		}
	}()

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

	if err := svr.ShutdownWithContext(ctx); err != nil {
		panic(err)
	}

	logs.Info("Server shutdown correctly")
}
