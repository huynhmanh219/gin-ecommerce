package config

import (
	"fmt"

	"gorm.io/gorm"
	"huynhmanh.com/gin/internal/model"
)

func RunMigrations(db *gorm.DB)error{
	if err:= db.AutoMigrate(
		&model.User{},
        &model.Product{},
        &model.Cart{},
        &model.CartItem{},
        &model.Order{},
        &model.OrderItem{},
        &model.Payment{},
        &model.InventoryLog{},
		); err != nil{
		return fmt.Errorf("Migration failed: %w",err)
	}
	fmt.Println("Migration completed")


	return nil
}