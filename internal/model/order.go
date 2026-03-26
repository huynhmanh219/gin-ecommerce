package model

import (
	"time"

	"github.com/shopspring/decimal"
)

const (
	OrderStatusPending    = "PENDING"
	OrderStatusProcessing = "PROCESSING"
	OrderStatusShipped    = "SHIPPED"
	OrderStatusDelivered  = "DELIVERED"
	OrderStatusCanceled   = "CANCELED"
)

type Order struct {
	ID              uint            `gorm:"primaryKey" json:"id"`
	OrderNumber     string          `gorm:"index" json:"order_number"`
	UserID          uint            `gorm:"index" json:"user_id"`
	TotalPrice      decimal.Decimal `gorm:"type:decimal(10,2)" json:"total_price"`
	Tax             decimal.Decimal `gorm:"type:decimal(10,2)" json:"tax"`
	ShippingFee     decimal.Decimal `gorm:"type:decimal(10,2)" json:"shipping_fee"`
	Discount        decimal.Decimal `gorm:"type:decimal(10,2)" json:"discount"`
	FinalPrice      decimal.Decimal `gorm:"type:decimal(10,2)" json:"final_price"`
	Status          string          `gorm:"index" json:"status"`
	ShippingAddress string          `json:"shipping_address"`
	Notes           string          `json:"notes"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`

	User    User        `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Items   []OrderItem `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"items"`
	Payment Payment     `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"payment,omitempty"`
}

func (Order) TableName() string {
	return "orders"
}

type OrderItem struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	OrderID   uint            `gorm:"index" json:"order_id"`
	ProductID uint            `json:"product_id"`
	Quantity  int             `json:"quantity"`
	UnitPrice decimal.Decimal `gorm:"type:decimal(10,2)" json:"unit_price"`
	Subtotal  decimal.Decimal `gorm:"type:decimal(10,2)" json:"sub_total"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`

	Order   *Order   `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"-"`
	Product Product `gorm:"foreignKey:ProductID;constraint:OnDelete:RESTRICT" json:"product,omitempty"`
}

func (OrderItem) TableName() string {
	return "order_items"
}
