package services

import (
	"errors"
	"testing"
	"time"

	"food_delivery/internal/configs"
	"food_delivery/internal/models/dto/requests"
	"food_delivery/internal/models/entities"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newOrderRepo() *mockOrderRepo {
	return &mockOrderRepo{
		createOrderWithTransactionFn: func(o *entities.Order) error { return nil },
		listByUserIDFn: func(uID uint, p, l int) ([]*entities.Order, int64, float64, map[string]int64, error) {
			return nil, 0, 0, nil, nil
		},
		findByOrderIDFn: func(id uint) (*entities.Order, error) { return nil, errors.New("not found") },
		updateFn:        func(o *entities.Order) error { return nil },
		listAllFn:       func(p, l int) ([]*entities.Order, int64, error) { return nil, 0, nil },
	}
}

func cartItemWithProduct(userID, productID uint, qty int, price float64) *entities.CartItem {
	item := &entities.CartItem{
		UserID:    userID,
		ProductID: productID,
		Quantity:  qty,
		Product:   entities.Product{Name: "Burger", Price: price},
	}
	return item
}

// ─── Checkout ─────────────────────────────────────────────────────────────────

func TestOrderService_Checkout_Success(t *testing.T) {
	orderRepo := newOrderRepo()
	cartRepo := newCartRepo()
	cartRepo.listByUserIDFn = func(uID uint) ([]*entities.CartItem, error) {
		return []*entities.CartItem{cartItemWithProduct(uID, 1, 2, 15.0)}, nil
	}
	svc := NewOrderService(orderRepo, cartRepo)

	resp, err := svc.Checkout(1)
	require.NoError(t, err)
	assert.Equal(t, 30.0, resp.TotalPrice)
	assert.Equal(t, "pending", resp.Status)
	assert.Len(t, resp.Items, 1)
}

func TestOrderService_Checkout_EmptyCart(t *testing.T) {
	orderRepo := newOrderRepo()
	cartRepo := newCartRepo()
	svc := NewOrderService(orderRepo, cartRepo)

	_, err := svc.Checkout(1)
	assert.Equal(t, configs.CartIsEmpty, err)
}

func TestOrderService_Checkout_CartError(t *testing.T) {
	orderRepo := newOrderRepo()
	cartRepo := newCartRepo()
	cartRepo.listByUserIDFn = func(uID uint) ([]*entities.CartItem, error) {
		return nil, errors.New("db error")
	}
	svc := NewOrderService(orderRepo, cartRepo)

	_, err := svc.Checkout(1)
	assert.Equal(t, configs.CartIsEmpty, err)
}

func TestOrderService_Checkout_CreateOrderFails(t *testing.T) {
	orderRepo := newOrderRepo()
	cartRepo := newCartRepo()
	cartRepo.listByUserIDFn = func(uID uint) ([]*entities.CartItem, error) {
		return []*entities.CartItem{cartItemWithProduct(uID, 1, 1, 10.0)}, nil
	}
	orderRepo.createOrderWithTransactionFn = func(o *entities.Order) error { return errors.New("db error") }
	svc := NewOrderService(orderRepo, cartRepo)

	_, err := svc.Checkout(1)
	assert.Equal(t, configs.CreateOrderFailed, err)
}

// ─── GetOrderHistory ──────────────────────────────────────────────────────────

func TestOrderService_GetOrderHistory_Success(t *testing.T) {
	orderRepo := newOrderRepo()
	cartRepo := newCartRepo()
	now := time.Now()
	orderRepo.listByUserIDFn = func(uID uint, p, l int) ([]*entities.Order, int64, float64, map[string]int64, error) {
		order := &entities.Order{TotalPrice: 50.0, Status: "completed"}
		order.CreatedAt = now
		order.OrderItems = []*entities.OrderItem{{ProductID: 1, Quantity: 1, Price: 50.0, Product: entities.Product{Name: "Burger"}}}
		return []*entities.Order{order}, 1, 50.0, map[string]int64{"completed": 1}, nil
	}
	svc := NewOrderService(orderRepo, cartRepo)

	resp, err := svc.GetOrderHistory(1, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), resp.TotalCount)
	assert.Equal(t, 50.0, resp.TotalPrice)
	assert.Equal(t, int64(1), resp.CompletedCount)
	assert.Equal(t, int64(0), resp.ProcessingCount)
}

func TestOrderService_GetOrderHistory_StatusCounts(t *testing.T) {
	orderRepo := newOrderRepo()
	cartRepo := newCartRepo()
	orderRepo.listByUserIDFn = func(uID uint, p, l int) ([]*entities.Order, int64, float64, map[string]int64, error) {
		return []*entities.Order{}, 0, 0, map[string]int64{
			"pending":    1,
			"processing": 2,
			"delivering": 3,
			"completed":  5,
		}, nil
	}
	svc := NewOrderService(orderRepo, cartRepo)

	resp, err := svc.GetOrderHistory(1, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(6), resp.ProcessingCount) // pending+processing+delivering
	assert.Equal(t, int64(5), resp.CompletedCount)
}

func TestOrderService_GetOrderHistory_Fails(t *testing.T) {
	orderRepo := newOrderRepo()
	cartRepo := newCartRepo()
	orderRepo.listByUserIDFn = func(uID uint, p, l int) ([]*entities.Order, int64, float64, map[string]int64, error) {
		return nil, 0, 0, nil, errors.New("db error")
	}
	svc := NewOrderService(orderRepo, cartRepo)

	_, err := svc.GetOrderHistory(1, 1, 10)
	assert.Equal(t, configs.FetchOrdersFailed, err)
}

// ─── AdminGetAllOrders ────────────────────────────────────────────────────────

func TestOrderService_AdminGetAllOrders_Success(t *testing.T) {
	orderRepo := newOrderRepo()
	cartRepo := newCartRepo()
	order := &entities.Order{TotalPrice: 20.0, Status: "pending"}
	order.OrderItems = []*entities.OrderItem{{ProductID: 1, Quantity: 1, Price: 20.0, Product: entities.Product{Name: "Pizza"}}}
	orderRepo.listAllFn = func(p, l int) ([]*entities.Order, int64, error) {
		return []*entities.Order{order}, 1, nil
	}
	svc := NewOrderService(orderRepo, cartRepo)

	resp, err := svc.AdminGetAllOrders(1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), resp.TotalCount)
	assert.Len(t, resp.Items, 1)
}

func TestOrderService_AdminGetAllOrders_Fails(t *testing.T) {
	orderRepo := newOrderRepo()
	cartRepo := newCartRepo()
	orderRepo.listAllFn = func(p, l int) ([]*entities.Order, int64, error) {
		return nil, 0, errors.New("db error")
	}
	svc := NewOrderService(orderRepo, cartRepo)

	_, err := svc.AdminGetAllOrders(1, 10)
	assert.Equal(t, configs.FetchOrdersFailed, err)
}

// ─── AdminUpdateOrderStatus ───────────────────────────────────────────────────

func TestOrderService_AdminUpdateOrderStatus_Success(t *testing.T) {
	orderRepo := newOrderRepo()
	cartRepo := newCartRepo()
	order := &entities.Order{Status: "pending"}
	order.ID = 1
	order.OrderItems = []*entities.OrderItem{}
	orderRepo.findByOrderIDFn = func(id uint) (*entities.Order, error) { return order, nil }
	svc := NewOrderService(orderRepo, cartRepo)

	resp, err := svc.AdminUpdateOrderStatus(1, requests.UpdateOrderStatusRequest{Status: "completed"})
	require.NoError(t, err)
	assert.Equal(t, "completed", resp.Status)
}

func TestOrderService_AdminUpdateOrderStatus_NotFound(t *testing.T) {
	orderRepo := newOrderRepo()
	cartRepo := newCartRepo()
	svc := NewOrderService(orderRepo, cartRepo)

	_, err := svc.AdminUpdateOrderStatus(9999, requests.UpdateOrderStatusRequest{Status: "completed"})
	assert.Equal(t, configs.OrderNotFound, err)
}

func TestOrderService_AdminUpdateOrderStatus_InvalidStatus(t *testing.T) {
	orderRepo := newOrderRepo()
	cartRepo := newCartRepo()
	order := &entities.Order{Status: "pending"}
	order.ID = 1
	orderRepo.findByOrderIDFn = func(id uint) (*entities.Order, error) { return order, nil }
	svc := NewOrderService(orderRepo, cartRepo)

	_, err := svc.AdminUpdateOrderStatus(1, requests.UpdateOrderStatusRequest{Status: "invalid_status"})
	assert.Equal(t, configs.InvalidOrderStatus, err)
}

func TestOrderService_AdminUpdateOrderStatus_UpdateFails(t *testing.T) {
	orderRepo := newOrderRepo()
	cartRepo := newCartRepo()
	order := &entities.Order{Status: "pending"}
	order.ID = 1
	orderRepo.findByOrderIDFn = func(id uint) (*entities.Order, error) { return order, nil }
	orderRepo.updateFn = func(o *entities.Order) error { return errors.New("db error") }
	svc := NewOrderService(orderRepo, cartRepo)

	_, err := svc.AdminUpdateOrderStatus(1, requests.UpdateOrderStatusRequest{Status: "completed"})
	assert.Equal(t, configs.UpdateOrderFailed, err)
}

func TestOrderService_AdminUpdateOrderStatus_AllValidStatuses(t *testing.T) {
	statuses := []string{"pending", "processing", "delivering", "completed", "cancelled"}
	for _, status := range statuses {
		t.Run(status, func(t *testing.T) {
			orderRepo := newOrderRepo()
			cartRepo := newCartRepo()
			order := &entities.Order{Status: "pending"}
			order.ID = 1
			order.OrderItems = []*entities.OrderItem{}
			orderRepo.findByOrderIDFn = func(id uint) (*entities.Order, error) { return order, nil }
			svc := NewOrderService(orderRepo, cartRepo)

			resp, err := svc.AdminUpdateOrderStatus(1, requests.UpdateOrderStatusRequest{Status: status})
			require.NoError(t, err)
			assert.Equal(t, status, resp.Status)
		})
	}
}
