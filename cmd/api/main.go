package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/oopbest/go-pos-backend/internal/config"
	"github.com/oopbest/go-pos-backend/internal/database"
)

func main() {
	// 1. โหลดการตั้งค่า
	cfg := config.LoadConfig()

	// 2. เชื่อมต่อฐานข้อมูล & Auto-migrate
	database.ConnectDB(cfg)

	// 3. สร้าง Fiber Server
	app := fiber.New()

	app.Use(logger.New())
	app.Use(cors.New())

	// Health Check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Restaurant POS Backend is connected to PostgreSQL!",
		})
	})

	log.Printf("Server starting on port %s...", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
