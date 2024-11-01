package controllers

import (
	"encoding/json"
	"io/ioutil"
	"os"
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
	filePath  = "data/userData.json"
)

func GetData(c *fiber.Ctx) error {
	mu.Lock()
	defer mu.Unlock()
	return c.JSON(dataStore)
}

func LoadData() {
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

func SaveData() {
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
