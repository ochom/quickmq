package app

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/basicauth"
	"github.com/ochom/gutils/env"
)

func New() *fiber.App {
	app := fiber.New()
	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("hello world")
	})

	auth := app.Use(basicauth.New(basicauth.Config{
		Users: map[string]string{
			env.Get("QUICK_MQ_USER", "admin"): env.Get("QUICK_MQ_PASSWORD", "admin"),
		},
	}))

	auth.Post("/publish", publisherHandler)
	auth.Get("/subscribe", subscriptionHandler)

	return app
}
