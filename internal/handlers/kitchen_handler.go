package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oopbest/go-pos-backend/internal/database"
	"github.com/oopbest/go-pos-backend/internal/models"
)

// @Summary      Get kitchen items
// @Description  ดึงรายการอาหารที่ต้องทำในครัว/บาร์ (pending, cooking, ready)
// @Tags         Kitchen
// @Produce      json
// @Param        station  query     string  false  "Station filter (kitchen or bar)"
// @Success      200      {array}   models.OrderItem
// @Router       /api/kitchen/items [get]
func GetKitchenItems(c *fiber.Ctx) error {
	station := c.Query("station")

	// ดึงเฉพาะจานที่ยังไม่เสิร์ฟ (pending, cooking, ready)
	query := database.DB.
		Preload("Product").
		Where("kitchen_status IN ?", []models.ItemKitchenStatus{
			models.ItemPending,
			models.ItemCooking,
			models.ItemReady,
		})

	if station != "" {
		query = query.Where("station = ?", station)
	}

	var items []models.OrderItem
	if err := query.Order("created_at asc").Find(&items).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch kitchen items",
		})
	}

	return c.JSON(items)
}

type UpdateKitchenStatusRequest struct {
	Status models.ItemKitchenStatus `json:"status"` // "pending", "cooking", "ready", "served"
}

// @Summary      Update item kitchen status
// @Description  อัปเดตสถานะของจานอาหาร (เช่น กำลังทำ, ทำเสร็จแล้ว, เสิร์ฟแล้ว)
// @Tags         Kitchen
// @Accept       json
// @Produce      json
// @Param        id    path      int                         true  "OrderItem ID"
// @Param        body  body      UpdateKitchenStatusRequest  true  "New Kitchen Status"
// @Success      200   {object}  models.OrderItem
// @Router       /api/kitchen/items/{id}/status [put]
func UpdateItemKitchenStatus(c *fiber.Ctx) error {
	itemID := c.Params("id")

	var req UpdateKitchenStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var item models.OrderItem
	if err := database.DB.First(&item, itemID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Order item not found"})
	}

	item.KitchenStatus = req.Status
	database.DB.Save(&item)

	database.DB.Preload("Product").First(&item, item.ID)
	return c.JSON(item)
}
