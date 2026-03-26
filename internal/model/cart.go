package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Cart struct {
	ID        uint 		`gorm:"primaryKey" json:"id"`
	UserID    uint 		`gorm:"index" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User User 			`gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Items []CartItem 	`gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE" json:"items"`
}

func (Cart) TableName() string{
	return "carts"
}

type CartItem struct{
	ID uint 					`gorm:"primaryKey" json:"id"`
	CartID uint					`gorm:"index" json:"cart_id"`
	ProductID uint 				`gorm:"index" json:"product_id"`
	Quantity int				`json:"quantity"`
	UnitPrice decimal.Decimal 	`gorm:"type:decimal(10,2)" json:"unit_price"`
	CreatedAt time.Time			`json:"created_at"`
	UpdatedAt time.Time			`json:"updated_at"`

	Cart Cart `gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE" json:"cart,omitempty"`
	Product Product `gorm:"foreignKey:ProductID;constraint:OnDelete:RESTRICT" json:"product,omitempty"`
}

func (CartItem) TableName() string {
	return "cart_items"
}