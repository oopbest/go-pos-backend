package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CategoryID  uint           `gorm:"not null" json:"category_id"`
	Category    Category       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Name        string         `gorm:"size:150;not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Price       float64        `gorm:"type:decimal(10,2);not null" json:"price"`
	Cost        float64        `gorm:"type:decimal(10,2);default:0.00" json:"cost"`
	ImageURL    string         `gorm:"size:255" json:"image_url"`
	IsAvailable bool           `gorm:"default:true" json:"is_available"` // สต็อกหมดหรือไม่
	Station     StationType    `gorm:"size:20;default:'kitchen'" json:"station"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
