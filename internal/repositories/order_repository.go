package repositories

import (
	"food_delivery/internal/models/entities"

	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrderWithTransaction(order *entities.Order) error
	ListByUserID(userID uint, page, limit int) ([]*entities.Order, int64, float64, map[string]int64, error)

	// Admin only
	FindByOrderID(id uint) (*entities.Order, error)
	Update(order *entities.Order) error
	ListAll(page, limit int) ([]*entities.Order, int64, error)
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

func (r *orderRepository) ListByUserID(userID uint, page, limit int) ([]*entities.Order, int64, float64, map[string]int64, error) {
	var orders []*entities.Order
	var totalCount int64
	var totalAmount float64

	if err := r.db.Model(&entities.Order{}).Where("user_id = ?", userID).Count(&totalCount).Error; err != nil {
		return nil, 0, 0, nil, err
	}

	if err := r.db.Model(&entities.Order{}).Where("user_id = ?", userID).Select("SUM(total_price)").Scan(&totalAmount).Error; err != nil {
		return nil, 0, 0, nil, err
	}

	var statusRows []struct {
		Status string `gorm:"column:status"`
		Count  int64  `gorm:"column:count"`
	}
	if err := r.db.Model(&entities.Order{}).
		Where("user_id = ?", userID).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&statusRows).Error; err != nil {
		return nil, 0, 0, nil, err
	}
	statusCounts := make(map[string]int64)
	for _, row := range statusRows {
		statusCounts[row.Status] = row.Count
	}

	if err := r.db.
		Preload("OrderItems").
		Preload("OrderItems.Product").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&orders).
		Error; err != nil {
		return nil, 0, 0, nil, err
	}

	return orders, totalCount, totalAmount, statusCounts, nil
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

func (r *orderRepository) ListAll(page, limit int) ([]*entities.Order, int64, error) {
	var orders []*entities.Order
	var totalCount int64
	if err := r.db.Model(&entities.Order{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.
		Preload("OrderItems").
		Preload("OrderItems.Product").
		Order("created_at DESC").
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&orders).
		Error; err != nil {
		return nil, 0, err
	}
	return orders, totalCount, nil
}
