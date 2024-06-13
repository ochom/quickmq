package app

import (
	"bufio"
	"fmt"
	"sync"

	"github.com/ochom/gutils/arrays"
	"github.com/ochom/gutils/helpers"
)

var channels map[string][]*client
var mu sync.Mutex

func init() {
	channels = make(map[string][]*client)
}

type client struct {
	id         string
	queue      string
	writer     *bufio.Writer
	tag        string
	stopSignal chan<- bool
}

func newClient(id string, writer *bufio.Writer, tag string, stopSignal chan<- bool) *client {
	return &client{id: id, writer: writer, tag: tag, stopSignal: stopSignal}
}

func (c *client) joinChannel(queueName string) {
	mu.Lock()
	defer mu.Unlock()
	c.queue = queueName
	channels[queueName] = append(channels[queueName], c)
}

func (c *client) leaveChannel() {
	mu.Lock()
	defer mu.Unlock()

	channels[c.queue] = arrays.Filter(channels[c.queue], func(cn *client) bool {
		return c.id != cn.id
	})
}

func (c *client) writeMessage(msg any) error {
	m := string(helpers.ToBytes(msg))
	_, err := c.writer.Write([]byte(fmt.Sprintf("data: %s\n\n", m)))
	if err != nil {
		return err
	}
	return c.writer.Flush()
}

func getMembers(queueName string) []*client {
	mu.Lock()
	defer mu.Unlock()

	return channels[queueName]
}
