package database

import (
	"fmt"
	"log"

	"github.com/oopbest/go-pos-backend/internal/config"
	"github.com/oopbest/go-pos-backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Bangkok",
		cfg.DBHost, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // แสดง SQL Query ใน Console
	})
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	log.Println(" Connected to PostgreSQL successfully!")

	// สั่ง Auto-Migrate สร้างตารางตาม Models
	err = db.AutoMigrate(
		&models.Table{},
		&models.Category{},
		&models.Product{},
		&models.Order{},
		&models.OrderItem{},
	)
	if err != nil {
		log.Fatalf("❌ AutoMigrate failed: %v", err)
	}
	log.Println(" Database Migration completed!")

	DB = db

	// Seed ข้อมูลเริ่มต้น (ถ้ายังไม่มีข้อมูล)
	SeedInitialData(db)

	return db
}

func SeedInitialData(db *gorm.DB) {
	var tableCount int64
	db.Model(&models.Table{}).Count(&tableCount)
	if tableCount == 0 {
		log.Println(" Seeding initial tables...")
		tables := []models.Table{
			{TableNo: "T-01", Zone: "Indoor", SeatCapacity: 4, Status: models.TableAvailable},
			{TableNo: "T-02", Zone: "Indoor", SeatCapacity: 4, Status: models.TableAvailable},
			{TableNo: "T-03", Zone: "Indoor", SeatCapacity: 2, Status: models.TableAvailable},
			{TableNo: "T-04", Zone: "Indoor", SeatCapacity: 6, Status: models.TableAvailable},
			{TableNo: "O-01", Zone: "Outdoor", SeatCapacity: 4, Status: models.TableAvailable},
			{TableNo: "O-02", Zone: "Outdoor", SeatCapacity: 4, Status: models.TableAvailable},
		}
		db.Create(&tables)
	}

	var catCount int64
	db.Model(&models.Category{}).Count(&catCount)
	if catCount == 0 {
		log.Println(" Seeding categories and products...")
		catFood := models.Category{Name: "อาหารจานเดียว", DefaultStation: models.StationKitchen}
		catDrink := models.Category{Name: "เครื่องดื่ม & กาแฟ", DefaultStation: models.StationBar}
		catDessert := models.Category{Name: "ของหวาน", DefaultStation: models.StationKitchen}
		db.Create(&catFood)
		db.Create(&catDrink)
		db.Create(&catDessert)

		products := []models.Product{
			{CategoryID: catFood.ID, Name: "ข้าวผัดกะเพราหมูกรอบ", Price: 79, Cost: 35, Station: models.StationKitchen, IsAvailable: true},
			{CategoryID: catFood.ID, Name: "ผัดไทยกุ้งสด", Price: 89, Cost: 40, Station: models.StationKitchen, IsAvailable: true},
			{CategoryID: catFood.ID, Name: "ข้าวไข่ข้นปู", Price: 120, Cost: 55, Station: models.StationKitchen, IsAvailable: true},
			{CategoryID: catDrink.ID, Name: "เอสเพรสโซ่เย็น", Price: 65, Cost: 20, Station: models.StationBar, IsAvailable: true},
			{CategoryID: catDrink.ID, Name: "มัทฉะลาเต้เย็น", Price: 75, Cost: 25, Station: models.StationBar, IsAvailable: true},
			{CategoryID: catDessert.ID, Name: "วาฟเฟิลไอศกรีม", Price: 85, Cost: 30, Station: models.StationKitchen, IsAvailable: true},
		}
		db.Create(&products)
	}
}
