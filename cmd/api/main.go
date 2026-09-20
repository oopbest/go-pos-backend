package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/oopbest/go-pos-backend/internal/config"
	"github.com/oopbest/go-pos-backend/internal/database"
	"github.com/oopbest/go-pos-backend/internal/handlers"
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

	// 4. API Routes Group
	api := app.Group("/api")

	// Table Routes
	api.Get("/tables", handlers.GetTables)
	api.Get("/tables/:id", handlers.GetTableByID)
	api.Put("/tables/:id/status", handlers.UpdateTableStatus)

	// Menu & Category Routes
	api.Get("/categories", handlers.GetCategories)
	api.Get("/products", handlers.GetProducts)

	log.Printf("Server starting on port %s...", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
