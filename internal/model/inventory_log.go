package model

import "time"

const (
	InventoryReasonOrder      = "ORDER"
	InventoryReasonAdjustment = "ADJUSTMENT"
	InventoryReasonReturn     = "RETURN"
	InventoryReasonRestock    = "RESTOCK"
)

type InventoryLog struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ProductID      uint      `gorm:"index" json:"product_id"`
	ChangeQuantity int       `json:"change_quantity"` // Dương = nhập, âm = xuất
	Reason         string    `json:"reason"`          // ORDER, ADJUSTMENT, RETURN, RESTOCK
	ReferenceID    uint      `json:"reference_id"`    // ID của order nếu reason = ORDER
	CreatedAt      time.Time `json:"created_at"`

	Product Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE" json:"product,omitempty"`
}


func (InventoryLog) TableName() string {
    return "inventory_logs"
}