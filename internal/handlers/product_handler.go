package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oopbest/go-pos-backend/internal/database"
	"github.com/oopbest/go-pos-backend/internal/models"
)

// @Summary      Get all categories
// @Description  ดึงหมวดหมู่ทั้งหมด
// @Tags         Menu
// @Produce      json
// @Success      200  {array}  models.Category
// @Router       /api/categories [get]
func GetCategories(c *fiber.Ctx) error {
	var categories []models.Category
	if err := database.DB.Order("id asc").Find(&categories).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch categories",
		})
	}
	return c.JSON(categories)
}

// @Summary      Get all products
// @Description  ดึงเมนูทั้งหมด (กรองตาม category_id หรือ station ได้)
// @Tags         Menu
// @Produce      json
// @Param        category_id  query     int     false  "Category ID"
// @Param        station      query     string  false  "Station (kitchen or bar)"
// @Success      200          {array}   models.Product
// @Router       /api/products [get]
func GetProducts(c *fiber.Ctx) error {
	categoryID := c.Query("category_id")
	station := c.Query("station")

	query := database.DB.Preload("Category").Where("is_available = ?", true)

	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if station != "" {
		query = query.Where("station = ?", station)
	}

	var products []models.Product
	if err := query.Order("id asc").Find(&products).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch products",
		})
	}
	return c.JSON(products)
}
