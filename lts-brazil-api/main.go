package main

import (
	"sync"

	"github.com/gofiber/fiber/v2"
)

type FormData struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Telephone string `json:"telephone"`
	Message   string `json:"message"`
}

var (
	dataStore []FormData
	mu        sync.Mutex
)

func main() {
	app := fiber.New()

	// Endpoint to receive data (POST)
	app.Post("/submit", func(c *fiber.Ctx) error {
		data := new(FormData)

		if err := c.BodyParser(data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot parse JSON",
			})
		}

		mu.Lock()
		dataStore = append(dataStore, *data)
		mu.Unlock()

		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Data received successfully",
			"data":    data,
		})
	})

	// Endpoint to list data (GET)
	app.Get("/data", func(c *fiber.Ctx) error {
		mu.Lock()
		defer mu.Unlock()
		return c.JSON(dataStore)
	})

	app.Listen(":3033")
}
