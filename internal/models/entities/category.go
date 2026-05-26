package entities

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name        string     `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
	Description string     `gorm:"type:text" json:"description"`
	Products    []*Product `gorm:"foreignKey:CategoryID" json:"-"`
}
