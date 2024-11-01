package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
	filePath  = "data/userData.json"
)

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowMethods: "GET,POST",
	}))

	loadData()

	app.Post("/api/users", func(c *fiber.Ctx) error {
		data := new(FormData)

		if err := c.BodyParser(data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot parse JSON",
			})
		}

		if data.Name == "" || data.Email == "" || data.Telephone == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Name, Email, and Telephone are required",
			})
		}

		mu.Lock()
		dataStore = append(dataStore, *data)
		saveData()
		mu.Unlock()

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"status":  "success",
			"message": "Data received successfully",
			"data":    data,
		})
	})

	app.Get("/data", func(c *fiber.Ctx) error {
		mu.Lock()
		defer mu.Unlock()
		return c.JSON(dataStore)
	})

	app.Get("/api/protected", func(c *fiber.Ctx) error {

		authorized := false // Change this according to your auth logic
		if !authorized {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized access",
			})
		}

		return c.JSON(fiber.Map{
			"message": "Welcome to the protected route!",
		})
	})

	app.Get("/api/restricted", func(c *fiber.Ctx) error {

		forbidden := true
		if forbidden {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access to this resource is forbidden",
			})
		}

		return c.JSON(fiber.Map{
			"message": "You have access!",
		})
	})

	app.Use(func(c *fiber.Ctx) error {
		if err := c.Next(); err != nil {
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal Server Error",
			})
			return err
		}
		return nil
	})

	app.Listen(":3033")
}

func loadData() {
	mu.Lock()
	defer mu.Unlock()

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return
	}

	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(file, &dataStore); err != nil {
		panic(err)
	}
}

func saveData() {

	if err := os.MkdirAll("data", os.ModePerm); err != nil {
		panic(err)
	}

	file, err := json.MarshalIndent(dataStore, "", "  ")
	if err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(filePath, file, 0644); err != nil {
		panic(err)
	}
}
