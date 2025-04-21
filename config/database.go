package config

import (
	"fmt"
	"log"
	"style-stamp/app/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDb() {
	dsn := "host=localhost user=postgres password=2034 dbname=stylestamp port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error!! , when connecting DB", err)
	}
	fmt.Println("DB Connected Successfully !!!")

	err = db.AutoMigrate(&models.Device{},&models.User{})
	if err != nil {
		log.Fatal("Error!! , when migrating DB", err)
	}
	DB = db
}
