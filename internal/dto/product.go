package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

type CreateProductRequest struct {
	Name        string `json:"name" validate:"required,min:2,max:100"`
	Description string `json:"description" validate:"min:2"`
	Price       decimal.Decimal `json:"price" validate:"required,gt=0"`
	Stock 		int `json:"stock"`
	Category 	string `json:"category"`
	ImageURL	string `json:"image_url"`
	IsActive 	bool	`json:"is_active" validate:"required"` 
}

type UpdateProductRequest struct{
	Name        string `json:"name"`
	Description string `json:"description" validate:"min:2"`
	Price       decimal.Decimal `json:"price" validate:"required,gt=0"`
	Stock 		int `json:"stock"`
	Category 	string `json:"category"`
	ImageURL	string `json:"image_url"`
	IsActive 	bool	`json:"is_active"` 
}

type ProductResponse struct{ 
	ID 			uint	`json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       decimal.Decimal `json:"price"`
	Stock 		int `json:"stock"`
	Category 	string `json:"category"`
	ImageURL	string `json:"image_url"`
	IsActive 	bool	`json:"is_active"` 
	CreatedAt	time.Time `json:"created_at"`
	UpdatedAt	time.Time `json:"updated_at"`
}
