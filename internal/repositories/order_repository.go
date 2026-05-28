package repositories

import (
	"food_delivery/internal/models/entities"

	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrderWithTransaction(order *entities.Order) error
	ListByUserID(userID uint) ([]*entities.Order, error)

	// Admin only
	FindByOrderID(id uint) (*entities.Order, error)
	Update(order *entities.Order) error
	ListAll() ([]*entities.Order, error)
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

func (r *orderRepository) ListByUserID(userID uint) ([]*entities.Order, error) {
	var orders []*entities.Order
	if err := r.db.
		Preload("OrderItems").
		Preload("OrderItems.Product").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&orders).
		Error; err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *orderRepository) FindByOrderID(id uint) (*entities.Order, error) {
	var order entities.Order
	if err := r.db.Preload("OrderItems").Preload("OrderItems.Product").First(&order, id).Error; err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *orderRepository) Update(order *entities.Order) error {
	return r.db.Save(order).Error
}

func (r *orderRepository) ListAll() ([]*entities.Order, error) {
	var orders []*entities.Order
	if err := r.db.
		Preload("OrderItems").
		Preload("OrderItems.Product").
		Order("created_at DESC").
		Find(&orders).
		Error; err != nil {
		return nil, err
	}
	return orders, nil
}
