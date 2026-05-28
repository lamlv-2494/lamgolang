package services

import (
	"food_delivery/internal/configs"
	"food_delivery/internal/models/dto/responses"
	"food_delivery/internal/models/entities"
	"food_delivery/internal/repositories"
)

type OrderService interface {
	Checkout(userID uint) (*responses.OrderResponse, error)
}

type orderService struct {
	orderRepo repositories.OrderRepository
	cartRepo  repositories.CartRepository
}

func NewOrderService(orderRepo repositories.OrderRepository, cartRepo repositories.CartRepository) OrderService {
	return &orderService{orderRepo: orderRepo, cartRepo: cartRepo}
}

func (s *orderService) Checkout(userID uint) (*responses.OrderResponse, error) {
	cartItems, err := s.cartRepo.ListByUserID(userID)
	if err != nil || len(cartItems) == 0 {
		return nil, configs.CartIsEmpty
	}

	var totalAmount float64
	var orderItems []*entities.OrderItem

	for _, item := range cartItems {
		priceSnapshot := item.Product.Price
		subTotal := priceSnapshot * float64(item.Quantity)
		totalAmount += subTotal

		orderItem := &entities.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     priceSnapshot,
		}
		orderItems = append(orderItems, orderItem)
	}

	order := &entities.Order{
		UserID:     userID,
		TotalPrice: totalAmount,
		Status:     "pending",
		OrderItems: orderItems,
	}

	if err := s.orderRepo.CreateOrderWithTransaction(order); err != nil {
		return nil, configs.CreateOrderFailed
	}

	var responseItems []responses.OrderItemData
	for i, item := range cartItems {
		savedOrderItem := order.OrderItems[i] // Lấy đúng ID vừa được DB cấp

		responseItems = append(responseItems, responses.OrderItemData{
			ID:        savedOrderItem.ID, // ID xịn của OrderItem
			ProductID: item.ProductID,
			Name:      item.Product.Name, // Tên lấy từ CartItem Preload
			Quantity:  savedOrderItem.Quantity,
			Price:     savedOrderItem.Price,
			SubTotal:  savedOrderItem.Price * float64(savedOrderItem.Quantity),
		})
	}

	return &responses.OrderResponse{
		ID:         order.ID,
		TotalPrice: order.TotalPrice,
		Status:     order.Status,
		CreatedAt:  order.CreatedAt,
		Items:      responseItems,
	}, nil
}
