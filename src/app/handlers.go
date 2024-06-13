package app

import (
	"bufio"

	"github.com/gofiber/fiber/v3"
	"github.com/ochom/gutils/uuid"
	"github.com/ochom/quickmq/src/domain"
)

func publisherHandler(c fiber.Ctx) error {
	var message domain.Message
	if err := c.Bind().Body(&message); err != nil {
		return err
	}

	if message.Queue == "" {
		return c.Status(400).JSON(fiber.Map{"error": "queue is required"})
	}

	if message.Body == nil {
		return c.Status(400).JSON(fiber.Map{"error": "body is required"})
	}

	// schedule message
	if message.Delay > 0 {
		scheduleMessage(message)
	}

	// publish instantly
	publishMessage(message)
	return c.JSON(fiber.Map{"status": "ok"})
}

func subscriptionHandler(c fiber.Ctx) error {
	queueName := c.Query("queue")
	if queueName == "" {
		return c.Status(400).JSON(fiber.Map{"error": "queue is required"})
	}

	id := uuid.New()
	tag := c.Query("tag", c.IP())

	stop := make(chan bool, 1)

	c.Context().SetContentType("text/event-stream")
	c.Context().Response.Header.Set("Cache-Control", "no-cache")
	c.Context().Response.Header.Set("Connection", "keep-alive")
	c.Context().Response.Header.Set("Transfer-Encoding", "chunked")
	c.Context().Response.Header.Set("Access-Control-Allow-Origin", "*")
	c.Context().Response.Header.Set("Access-Control-Allow-Headers", "Cache-Control")
	c.Context().Response.Header.Set("Access-Control-Allow-Credentials", "true")

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		cl := newClient(id, w, tag, stop)
		go func() {
			cl.joinChannel(queueName)
		}()

		cl.writeMessage(fiber.Map{"status": "connected"})
		<-stop
	})

	return nil
}
