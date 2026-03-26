package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Product struct {
	ID          uint			`gorm:"primaryKey" json:"id"`
	Name        string			`gorm:"index" json:"name"`
	Description string			`json:"description"`
	Price       decimal.Decimal `gorm:"type:decimal(10,2)" json:"price"`
	Stock       int				`json:"stock"`
	Category string				`gorm:"index" json:"category"`
	ImageURL string				`json:"image_url"`
	IsActive bool				`gorm:"default:true" json:"is_active"`
	CreatedAt time.Time			`json:"created_at"`
	UpdatedAt time.Time			`json:"updated_at"`

	CartItems []CartItem 		`gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE" json:"-"`
	OrderItems []OrderItem		`gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE" json:"-"`
	InventoryLogs []InventoryLog `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE" json:"-"`
}

func (Product) TableName() string{
	return "products"
}