package api

import (
	"github.com/gofiber/fiber/v3"
)

func New() *fiber.App {
	app := fiber.New()

	// rest apis
	app.Post("/login", login)
	app.Get("/user", loadUSer)

	// serve other static files
	app.Static("/", "web/build")
	app.Get("*", func(c fiber.Ctx) error {
		return c.SendFile("web/build/index.html")
	})

	return app
}
