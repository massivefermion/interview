package database

import (
	"interview/entities"

	"gorm.io/gorm"
)

func MigrateDatabase(db *gorm.DB) {
	// AutoMigrate will create or update the tables based on the models
	err := db.AutoMigrate(&entities.CartEntity{}, &entities.CartItem{})
	if err != nil {
		panic(err)
	}
}
