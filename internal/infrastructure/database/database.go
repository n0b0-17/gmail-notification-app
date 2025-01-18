package database

import (
  "gmail-notification-app/internal/domain/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


func NewDatabase(dsn string) (*gorm.DB, error){
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err 
	}

	//自動マイグレートスキーマ
	err = db.AutoMigrate(&models.User{})
	if err != nil{
		log.Printf("Failed to migrate database: %v", err)
		return nil, err
	}

	return db,nil
}


