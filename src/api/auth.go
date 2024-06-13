package api

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/ochom/gutils/auth"
	"github.com/ochom/gutils/env"
	"github.com/ochom/gutils/uuid"
)

func login(c fiber.Ctx) error {
	var data struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.Bind().Body(&data); err != nil {
		return err
	}

	username := env.Get("QUICK_MQ_USER", "admin")
	password := env.Get("QUICK_MQ_PASSWORD", "admin")

	if data.Username != username {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid username"})
	}

	if data.Password != password {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid password"})
	}

	token, err := auth.GenerateAuthTokens(map[string]string{"user": "admin", "session_id": uuid.New()})
	if err != nil {
		return err
	}

	// Set a cookie
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token["token"],
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	return c.JSON(fiber.Map{"status": "success", "message": "Logged in", "username": "admin"})
}

func loadUSer(c fiber.Ctx) error {
	token := c.Cookies("jwt")
	claims, err := auth.GetAuthClaims(token)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "Unauthorized", "data": err.Error()})
	}

	if claims["user"] != "admin" {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "Unauthorized"})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Logged in", "username": "admin"})
}
