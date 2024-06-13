package app

import (
	"time"

	"github.com/ochom/gutils/cache"
	"github.com/ochom/gutils/logs"
	"github.com/ochom/gutils/uuid"
	"github.com/ochom/quickmq/src/domain"
)

// Consume runs the consumers in the background
func Consume(stop <-chan bool) {
	for {
		select {
		case msg := <-readMessages():
			handleMessage(msg)
		case <-stop:
			return
		}
	}
}

func handleMessage(msg domain.Message) {
	cls := getMembers(msg.Queue)
	// if there are no clients, reschedule message to be sent after 5 seconds
	if len(cls) == 0 {
		cache.SetWithCallback(uuid.New(), []byte("reschedule"), time.Second*5, func() {
			publishMessage(msg)
		})

		return
	}

	for _, c := range cls {
		if err := c.writeMessage(msg.Body); err != nil {
			logs.Error("Error writing message: %s", err.Error())
			c.leaveChannel()
			return
		}
	}
}
