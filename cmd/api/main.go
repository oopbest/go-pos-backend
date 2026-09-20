package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New()) // เพื่อให้ React Frontend เรียก API ข้าม Port ได้

	// Health Check Route
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Restaurant POS Backend is running!",
		})
	})

	log.Println("Server starting on port 8080...")
	log.Fatal(app.Listen(":8080"))
}
