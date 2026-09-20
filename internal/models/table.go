package models

import (
	"time"

	"gorm.io/gorm"
)

type TableStatus string

const (
	TableAvailable TableStatus = "available" // โต๊ะว่าง
	TableOccupied  TableStatus = "occupied"  // มีลูกค้านั่ง
	TableReserved  TableStatus = "reserved"  // จองไว้
)

type Table struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	TableNo      string         `gorm:"size:20;uniqueIndex;not null" json:"table_no"` // e.g. "T-01", "A1"
	Zone         string         `gorm:"size:50;default:'Indoor'" json:"zone"`         // e.g. "Indoor", "Outdoor", "Bar"
	SeatCapacity int            `gorm:"default:4" json:"seat_capacity"`
	Status       TableStatus    `gorm:"size:20;default:'available'" json:"status"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
