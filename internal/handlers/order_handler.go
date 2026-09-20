package handlers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/oopbest/go-pos-backend/internal/database"
	"github.com/oopbest/go-pos-backend/internal/models"
)

// DTO สำหรับรับข้อมูลสั่งอาหาร
type CreateOrderItemRequest struct {
	ProductID           uint   `json:"product_id"`
	Quantity            int    `json:"quantity"`
	SpecialInstructions string `json:"special_instructions"`
}

type CreateOrderRequest struct {
	TableID       uint                     `json:"table_id"`
	CustomerCount int                      `json:"customer_count"`
	Items         []CreateOrderItemRequest `json:"items"`
}

// @Summary      Create order / Open table
// @Description  เปิดโต๊ะและสั่งรายการอาหารครั้งแรก
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        body  body      CreateOrderRequest  true  "Order Payload"
// @Success      201   {object}  models.Order
// @Router       /api/orders [post]
func CreateOrder(c *fiber.Ctx) error {
	var req CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if len(req.Items) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Order must contain at least one item"})
	}

	// เริ่ม Database Transaction
	tx := database.DB.Begin()

	// 1. ตรวจสอบว่าโต๊ะว่างหรือไม่
	var table models.Table
	if err := tx.First(&table, req.TableID).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Table not found"})
	}

	if table.Status == models.TableOccupied {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Table is already occupied"})
	}

	// 2. สร้าง Order No. แบบไม่ซ้ำ เช่น ORD-20260920-150405
	orderNo := fmt.Sprintf("ORD-%s", time.Now().Format("20060102-150405"))

	order := models.Order{
		OrderNo:       orderNo,
		TableID:       req.TableID,
		CustomerCount: req.CustomerCount,
		Status:        models.OrderOpen,
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create order"})
	}

	// 3. เพิ่มรายการอาหาร และคำนวณราคารวม
	var subtotal float64
	for _, itemReq := range req.Items {
		var product models.Product
		if err := tx.First(&product, itemReq.ProductID).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Product ID %d not found", itemReq.ProductID)})
		}

		itemSubtotal := product.Price * float64(itemReq.Quantity)
		subtotal += itemSubtotal

		orderItem := models.OrderItem{
			OrderID:             order.ID,
			ProductID:           product.ID,
			Quantity:            itemReq.Quantity,
			UnitPrice:           product.Price,
			Subtotal:            itemSubtotal,
			SpecialInstructions: itemReq.SpecialInstructions,
			KitchenStatus:       models.ItemPending,
			Station:             product.Station,
		}

		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create order item"})
		}
	}

	// 4. อัปเดตยอดเงินรวมของ Order
	order.Subtotal = subtotal
	order.TotalAmount = subtotal // ในขั้นตอนนี้ยอดรวม = subtotal (ยังไม่หัก discount / vat)
	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update order total"})
	}

	// 5. เปลี่ยนสถานะโต๊ะเป็น "occupied" (มีลูกค้า)
	table.Status = models.TableOccupied
	if err := tx.Save(&table).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update table status"})
	}

	// Commit Transaction
	tx.Commit()

	// โหลดข้อมูล Order พร้อม Items และ Table ส่งกลับไป
	database.DB.Preload("Table").Preload("Items.Product").First(&order, order.ID)

	return c.Status(fiber.StatusCreated).JSON(order)
}

// @Summary      Get active order by table
// @Description  ดึงบิลปัจจุบันที่ยังเปิดอยู่ของโต๊ะนั้น
// @Tags         Orders
// @Produce      json
// @Param        table_id  path      int  true  "Table ID"
// @Success      200       {object}  models.Order
// @Router       /api/orders/table/{table_id} [get]
func GetActiveOrderByTableID(c *fiber.Ctx) error {
	tableID := c.Params("table_id")

	var order models.Order
	err := database.DB.
		Preload("Table").
		Preload("Items.Product").
		Where("table_id = ? AND status = ?", tableID, models.OrderOpen).
		First(&order).Error

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No active order for this table",
		})
	}

	return c.JSON(order)
}
