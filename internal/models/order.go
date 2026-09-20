package models

import (
	"time"

	"gorm.io/gorm"
)

type OrderStatus string
type ItemKitchenStatus string

const (
	OrderOpen      OrderStatus = "open"      // กำลังทาน ยังไม่เช็คบิล
	OrderCompleted OrderStatus = "completed" // เช็คบิลแล้ว
	OrderCancelled OrderStatus = "cancelled" // ยกเลิก

	ItemPending ItemKitchenStatus = "pending" // เพิ่งสั่ง ยังไม่เริ่มทำ
	ItemCooking ItemKitchenStatus = "cooking" // ครัวกำลังปรุง
	ItemReady   ItemKitchenStatus = "ready"   // ทำเสร็จ รอเสิร์ฟ
	ItemServed  ItemKitchenStatus = "served"  // เสิร์ฟแล้ว
)

type Order struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	OrderNo       string         `gorm:"size:50;uniqueIndex;not null" json:"order_no"` // e.g. "ORD-20260920-001"
	TableID       uint           `gorm:"not null" json:"table_id"`
	Table         Table          `gorm:"foreignKey:TableID" json:"table,omitempty"`
	CustomerCount int            `gorm:"default:1" json:"customer_count"`
	Status        OrderStatus    `gorm:"size:20;default:'open'" json:"status"`
	Subtotal      float64        `gorm:"type:decimal(10,2);default:0.00" json:"subtotal"`
	Discount      float64        `gorm:"type:decimal(10,2);default:0.00" json:"discount"`
	Tax           float64        `gorm:"type:decimal(10,2);default:0.00" json:"tax"`
	TotalAmount   float64        `gorm:"type:decimal(10,2);default:0.00" json:"total_amount"`
	PaymentMethod string         `gorm:"size:30" json:"payment_method"` // "cash", "promptpay", "credit_card"
	PaidAt        *time.Time     `json:"paid_at"`
	Items         []OrderItem    `gorm:"foreignKey:OrderID" json:"items"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type OrderItem struct {
	ID                  uint              `gorm:"primaryKey" json:"id"`
	OrderID             uint              `gorm:"not null;index" json:"order_id"`
	ProductID           uint              `gorm:"not null" json:"product_id"`
	Product             Product           `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Quantity            int               `gorm:"not null;default:1" json:"quantity"`
	UnitPrice           float64           `gorm:"type:decimal(10,2);not null" json:"unit_price"`
	Subtotal            float64           `gorm:"type:decimal(10,2);not null" json:"subtotal"`
	SpecialInstructions string            `gorm:"size:255" json:"special_instructions"` // e.g. "ไม่หวาน, แยกน้ำแข็ง"
	KitchenStatus       ItemKitchenStatus `gorm:"size:20;default:'pending'" json:"kitchen_status"`
	Station             StationType       `gorm:"size:20;default:'kitchen'" json:"station"`
	CreatedAt           time.Time         `json:"created_at"`
	UpdatedAt           time.Time         `json:"updated_at"`
}
