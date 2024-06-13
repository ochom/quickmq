package app

import (
	"time"

	"github.com/ochom/gutils/cache"
	"github.com/ochom/gutils/uuid"
	"github.com/ochom/quickmq/src/domain"
)

var messages = make(chan domain.Message, 1000)

func publishMessage(msg domain.Message) {
	messages <- msg
}

func scheduleMessage(msg domain.Message) {
	cache.SetWithCallback(uuid.New(), []byte("delayed-message"), time.Second*time.Duration(msg.Delay), func() {
		messages <- msg
	})
}

func readMessages() chan domain.Message {
	return messages
}
