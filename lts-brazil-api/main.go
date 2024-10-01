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
	filePath  = "data/formData.json"
)

func main() {
	app := fiber.New()

	// Adicionar middleware de CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowMethods: "GET,POST",
	}))

	// Carregar dados do arquivo JSON ao iniciar
	loadData()

	// Endpoint para receber dados (POST)
	app.Post("/submit", func(c *fiber.Ctx) error {
		data := new(FormData)

		if err := c.BodyParser(data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot parse JSON",
			})
		}

		mu.Lock()
		dataStore = append(dataStore, *data)
		saveData()
		mu.Unlock()

		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Data received successfully",
			"data":    data,
		})
	})

	// Endpoint para listar dados (GET)
	app.Get("/data", func(c *fiber.Ctx) error {
		mu.Lock()
		defer mu.Unlock()
		return c.JSON(dataStore)
	})

	app.Listen(":3033")
}

// Função para carregar dados do arquivo JSON
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

// Função para salvar dados no arquivo JSON
func saveData() {
	// Ensure the directory exists
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
