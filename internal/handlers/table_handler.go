package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oopbest/go-pos-backend/internal/database"
	"github.com/oopbest/go-pos-backend/internal/models"
)

// @Summary      Get all tables
// @Description  ดึงรายชื่อโต๊ะทั้งหมด
// @Tags         Tables
// @Produce      json
// @Success      200  {array}   models.Table
// @Router       /api/tables [get]
func GetTables(c *fiber.Ctx) error {
	var tables []models.Table
	if err := database.DB.Order("table_no asc").Find(&tables).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch tables",
		})
	}
	return c.JSON(tables)
}

// @Summary      Get table by ID
// @Description  ดึงข้อมูลโต๊ะเดี่ยวๆ ตาม ID
// @Tags         Tables
// @Produce      json
// @Param        id   path      int  true  "Table ID"
// @Success      200  {object}  models.Table
// @Router       /api/tables/{id} [get]
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

// @Summary      Update table status
// @Description  อัปเดตสถานะโต๊ะ (available, occupied, reserved)
// @Tags         Tables
// @Accept       json
// @Produce      json
// @Param        id    path      int     true  "Table ID"
// @Param        body  body      object  true  "Status Payload (e.g. {\"status\": \"available\"})"
// @Success      200   {object}  models.Table
// @Router       /api/tables/{id}/status [put]
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
