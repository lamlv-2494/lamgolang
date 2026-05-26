package entities

import "time"

type Rating struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
	ProductID uint      `gorm:"not null" json:"product_id"`
	Stars     int       `gorm:"not null" json:"stars"` // Từ 1 -> 5 sao
	Comment   string    `gorm:"type:text" json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
