package main

import (
	"lts-brazil-api/controllers"
	"lts-brazil-api/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowMethods: "GET,POST",
	}))

	userService := services.NewUserService()

	app.Post("/api/users", controllers.PostUser(userService))
	app.Get("/data", controllers.GetData)
	app.Get("/api/protected", controllers.GetProtected)
	app.Get("/api/restricted", controllers.GetRestricted)

	app.Use(func(c *fiber.Ctx) error {
		if err := c.Next(); err != nil {
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
			return err
		}
		return nil
	})

	app.Listen(":3033")
}
