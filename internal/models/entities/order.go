package entities

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserID     uint         `json:"user_id"`
	User       User         `gorm:"foreignKey:UserID" json:"-"`
	TotalPrice float64      `gorm:"type:decimal(10,2);not null" json:"total_price"`
	Status     string       `gorm:"type:varchar(50);default:'pending'" json:"status"` // pending, processing, completed, cancelled
	OrderItems []*OrderItem `gorm:"foreignKey:OrderID" json:"order_items"`
}
