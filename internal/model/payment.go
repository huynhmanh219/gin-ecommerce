package model

import (
	"time"

	"github.com/shopspring/decimal"
	datatypes "gorm.io/datatypes"
)

const (
	PaymentStatusPending   = "PENDING"
	PaymentStatusCompleted = "COMPLETED"
	PaymentStatusFailed    = "FAILED"
	PaymentStatusRefunded  = "REFUNDED"
)

const (
	PaymentMethodCreditCard   = "CREDIT_CARD"
	PaymentMethodBankTransfer = "BANK_TRANSFER"
	PaymentMethodCash         = "CASH"
)

type Payment struct {
	ID              uint            `gorm:"primaryKey" json:"id"`
	OrderID         uint            `gorm:"index" json:"order_id"` // 1-1 với Order
	Amount          decimal.Decimal `gorm:"type:decimal(10,2)" json:"amount"`
	Method          string          `json:"method"`              // CREDIT_CARD, BANK_TRANSFER, CASH
	Status          string          `gorm:"index" json:"status"` // PENDING, COMPLETED, FAILED, REFUNDED
	TransactionID   string          `json:"transaction_id"`      // ID từ payment gateway
	GatewayResponse datatypes.JSON  `json:"gateway_response"`    // Response từ payment provider (JSON)
	PaidAt          *time.Time      `json:"paid_at"`             // NULL nếu chưa thanh toán
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`

	// Relationships
	Order *Order `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"-"`
}

func (Payment) TableName() string {
	return "payments"
}
