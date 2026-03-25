package config

import (
	"fmt"

	"gorm.io/gorm"
	"huynhmanh.com/gin/internal/model"
)

func RunMigrations(db *gorm.DB)error{
	if err:= db.AutoMigrate(&model.User{}); err != nil{
		return fmt.Errorf("Migration failed: %w",err)
	}
	if err:= db.AutoMigrate(&model.AuditLog{}); err != nil{
		return fmt.Errorf("Migration failed: %w",err)
	}
	if err:= db.AutoMigrate(&model.UserProFile{}); err != nil{
		return fmt.Errorf("Migration failed: %w",err)
	}
	fmt.Println("Migration completed")


	return nil
}