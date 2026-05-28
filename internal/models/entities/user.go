package entities

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(100);uniqueIndex;not null" json:"username" binding:"required"`
	Email    string `gorm:"type:varchar(100);uniqueIndex;not null" json:"email" binding:"required,email"`
	Password string `gorm:"not null" json:"-" binding:"required"`
	Role     string `gorm:"type:varchar(20);default:'user'" json:"role"` // 'user' or 'admin'

	CartItems   []*CartItem   `gorm:"foreignKey:UserID" json:"-"`
	Orders      []*Order      `gorm:"foreignKey:UserID" json:"-"`
	Ratings     []*Rating     `gorm:"foreignKey:UserID" json:"-"`
	Suggestions []*Suggestion `gorm:"foreignKey:UserID" json:"-"`
}
