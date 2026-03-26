package model

import "time"

type User struct {
	ID        uint   	`gorm:"primaryKey" json:"id"`
	Name      string 	`gorm:"index" json:"name"`
	Email     string 	`gorm:"uniqueIndex;type:varchar(255)" json:"email"`
	Password  string 	`json:"-"`
	IsActive  bool   	`gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time	`json:"updated_at"`

	Cart    *Cart   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"cart,omitempty"`
	Orders  []Order `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"orders,omitempty"`
}

type UserProFile struct{
	ID uint `gorm:"primaryKey"`
	UserID uint `gorm:"uniqueIndex"`
	Bio string `json:"bio"`
	Avatar string `json:"avatar"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AuditLog struct {
	ID uint `gorm:"primaryKey"`
	UserID uint
	Action string
	Details string
	IPAddress string
	CreatedAt time.Time
}



func (User) TableName() string{
	return "users"
}