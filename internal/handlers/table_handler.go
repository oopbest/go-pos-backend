package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oopbest/go-pos-backend/internal/database"
	"github.com/oopbest/go-pos-backend/internal/models"
)

// GET /api/tables - ดึงรายชื่อโต๊ะทั้งหมด เรียงตามเลขโต๊ะ
func GetTables(c *fiber.Ctx) error {
	var tables []models.Table
	if err := database.DB.Order("table_no asc").Find(&tables).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch tables",
		})
	}
	return c.JSON(tables)
}

// GET /api/tables/:id - ดึงข้อมูลโต๊ะเดี่ยวๆ ตาม ID
func GetTableByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var table models.Table
	if err := database.DB.First(&table, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Table not found",
		})
	}
	return c.JSON(table)
}

// PUT /api/tables/:id/status - อัปเดตสถานะโต๊ะ (available, occupied, reserved)
func UpdateTableStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	type StatusRequest struct {
		Status models.TableStatus `json:"status"`
	}

	var req StatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	var table models.Table
	if err := database.DB.First(&table, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Table not found",
		})
	}

	table.Status = req.Status
	database.DB.Save(&table)

	return c.JSON(table)
}
