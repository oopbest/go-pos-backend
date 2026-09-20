package models

import (
	"time"

	"gorm.io/gorm"
)

type StationType string

const (
	StationKitchen StationType = "kitchen" // ครัวอาหาร
	StationBar     StationType = "bar"     // บาร์น้ำ/กาแฟ
)

type Category struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Name           string         `gorm:"size:100;not null" json:"name"` // e.g. "อาหารจานเดียว", "กาแฟ", "ของหวาน"
	DefaultStation StationType    `gorm:"size:20;default:'kitchen'" json:"default_station"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}
