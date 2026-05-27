package entities

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string       `gorm:"type:varchar(255);not null" json:"name"`
	Description string       `gorm:"type:text" json:"description"`
	Price       float64      `gorm:"type:decimal(10,2);not null" json:"price"`
	Image       string       `gorm:"type:varchar(255)" json:"image"`
	CategoryID  uint         `json:"category_id"`
	Category    Category     `gorm:"foreignKey:CategoryID" json:"category"`
	CartItems   []*CartItem  `gorm:"foreignKey:ProductID" json:"-"`
	OrderItems  []*OrderItem `gorm:"foreignKey:ProductID" json:"-"`
	Ratings     []*Rating    `gorm:"foreignKey:ProductID" json:"ratings,omitempty"`
}
