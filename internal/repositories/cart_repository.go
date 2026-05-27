package repositories

import (
	"food_delivery/internal/models/entities"

	"gorm.io/gorm"
)

type CartRepository interface {
	Create(item *entities.CartItem) error
	FindByID(id uint) (*entities.CartItem, error)
	FindByUserAndProduct(userID uint, productID uint) (*entities.CartItem, error)
	Update(item *entities.CartItem) error
	Delete(item *entities.CartItem) error
	ListByUserID(userID uint) ([]*entities.CartItem, error)
}

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) Create(item *entities.CartItem) error {
	return r.db.Create(item).Error
}

func (r *cartRepository) FindByID(id uint) (*entities.CartItem, error) {
	var item entities.CartItem
	if err := r.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *cartRepository) FindByUserAndProduct(userID uint, productID uint) (*entities.CartItem, error) {
	var item entities.CartItem
	if err := r.db.Where("user_id = ? AND product_id = ?", userID, productID).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *cartRepository) Update(item *entities.CartItem) error {
	return r.db.Save(item).Error
}

func (r *cartRepository) Delete(item *entities.CartItem) error {
	return r.db.Delete(item).Error
}

func (r *cartRepository) ListByUserID(userID uint) ([]*entities.CartItem, error) {
	var items []*entities.CartItem
	if err := r.db.Where("user_id = ?", userID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
