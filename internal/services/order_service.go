package services

import (
	"food_delivery/internal/configs"
	"food_delivery/internal/models/dto/requests"
	"food_delivery/internal/models/dto/responses"
	"food_delivery/internal/models/entities"
	"food_delivery/internal/repositories"
)

type OrderService interface {
	Checkout(userID uint) (*responses.OrderResponse, error)

	GetOrderHistory(userID uint) ([]*responses.OrderResponse, error)

	AdminGetAllOrders() ([]*responses.OrderResponse, error)
	AdminUpdateOrderStatus(orderID uint, req requests.UpdateOrderStatusRequest) (*responses.OrderResponse, error)
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
		savedOrderItem := order.OrderItems[i]

		responseItems = append(responseItems, responses.OrderItemData{
			ID:        savedOrderItem.ID,
			ProductID: item.ProductID,
			Name:      item.Product.Name,
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

func (s *orderService) GetOrderHistory(userID uint) ([]*responses.OrderResponse, error) {
	orders, err := s.orderRepo.ListByUserID(userID)
	if err != nil {
		return nil, configs.FetchOrdersFailed
	}

	return s.mapOrdersToResponse(orders), nil
}

func (s *orderService) AdminGetAllOrders() ([]*responses.OrderResponse, error) {
	orders, err := s.orderRepo.ListAll()
	if err != nil {
		return nil, configs.FetchOrdersFailed
	}

	return s.mapOrdersToResponse(orders), nil
}

func (s *orderService) mapOrdersToResponse(orders []*entities.Order) []*responses.OrderResponse {
	var result []*responses.OrderResponse = make([]*responses.OrderResponse, 0)
	for _, order := range orders {
		var responseItems []responses.OrderItemData
		for _, item := range order.OrderItems {
			responseItems = append(responseItems, responses.OrderItemData{
				ID:        item.ID,
				ProductID: item.ProductID,
				Name:      item.Product.Name,
				Quantity:  item.Quantity,
				Price:     item.Price,
				SubTotal:  item.Price * float64(item.Quantity),
			})
		}
		result = append(result, &responses.OrderResponse{
			ID:         order.ID,
			TotalPrice: order.TotalPrice,
			Status:     order.Status,
			CreatedAt:  order.CreatedAt,
			Items:      responseItems,
		})
	}
	return result
}

func (s *orderService) AdminUpdateOrderStatus(orderID uint, req requests.UpdateOrderStatusRequest) (*responses.OrderResponse, error) {
	order, err := s.orderRepo.FindByOrderID(orderID)
	if err != nil {
		return nil, configs.OrderNotFound
	}

	// only accept valid status values
	validStatuses := map[string]bool{
		"pending": true, "processing": true, "delivering": true, "completed": true, "cancelled": true,
	}
	if !validStatuses[req.Status] {
		return nil, configs.InvalidOrderStatus
	}

	order.Status = req.Status
	if err := s.orderRepo.Update(order); err != nil {
		return nil, configs.UpdateOrderFailed
	}

	return s.mapOrdersToResponse([]*entities.Order{order})[0], nil
}
