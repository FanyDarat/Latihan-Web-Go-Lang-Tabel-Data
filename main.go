package main

import "github.com/gofiber/fiber/v2"

func main() {
	app := fiber.New()

	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello, Fiber!",
		})
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to the Fiber server!")
	})

	app.Listen(":8080")
}
