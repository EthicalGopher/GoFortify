package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("index.html")
	})

	app.Get("/api/message", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello from GoFortify example backend!",
		})
	})

	app.Post("/api/echo", func(c *fiber.Ctx) error {
		type Request struct {
			Name     string `json:"name"`
			Password string `json:"password"`
		}
		var req Request
		if err := c.BodyParser(&req); err != nil {
			return err
		}
		return c.JSON(fiber.Map{
			"name":     req.Name,
			"password": req.Password,
		})
	})

	app.Get("/api/search", func(c *fiber.Ctx) error {
		q := c.Query("q")
		return c.JSON(fiber.Map{
			"query":   q,
			"results": []string{"Result 1", "Result 2"},
		})
	})

	log.Println("Example backend server started on :3001")
	log.Fatal(app.Listen(":3001"))
}
