package controllers

import (
	"lts-brazil-api/models"
	"lts-brazil-api/services"

	"github.com/gofiber/fiber/v2"
)

type User struct {
	FullName       string `json:"fullName"`
	CPF            string `json:"cpf"`
	BirthDate      string `json:"birthDate"`
	PhoneNumber    string `json:"phoneNumber"`
	RoleName       string `json:"roleName"`
	PermissionType string `json:"permissionType"`
}

func PostUser(userService *services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := new(User)
		if err := c.BodyParser(user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}
		if user.FullName == "" || user.CPF == "" || user.BirthDate == "" || user.PhoneNumber == "" || user.RoleName == "" || user.PermissionType == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "All fields are required"})
		}

		if err := models.ValidateCPF(user.CPF); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		err := userService.CreateUser(user.FullName, user.CPF, user.BirthDate, user.PhoneNumber, user.RoleName, user.PermissionType)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "message": "User created successfully", "data": user})
	}
}

func GetProtected(c *fiber.Ctx) error {
	authorized := false // Change this according to your auth logic
	if !authorized {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized access"})
	}
	return c.JSON(fiber.Map{"message": "Welcome to the protected route!"})
}

func GetRestricted(c *fiber.Ctx) error {
	forbidden := true
	if forbidden {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access to this resource is forbidden"})
	}
	return c.JSON(fiber.Map{"message": "You have access!"})
}
