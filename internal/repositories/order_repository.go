package repositories

import (
	"food_delivery/internal/models/entities"

	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrderWithTransaction(order *entities.Order) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateOrderWithTransaction(order *entities.Order) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		if err := tx.Where("user_id = ?", order.UserID).Delete(&entities.CartItem{}).Error; err != nil {
			return err
		}

		return nil
	})
}
