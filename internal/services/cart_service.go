package services

import (
	"food_delivery/internal/configs"
	"food_delivery/internal/models/dto/requests"
	"food_delivery/internal/models/dto/responses"
	"food_delivery/internal/models/entities"
	"food_delivery/internal/repositories"
)

type CartService interface {
	AddToCart(userID uint, req *requests.AddToCartRequest) error
	GetCart(userID uint) (*responses.CartResponse, error)
	UpdateCartItem(userID uint, cartItemID uint, req *requests.UpdateCartItemRequest) error
	RemoveFromCart(userID uint, cartItemID uint) error
}

type cartService struct {
	cartRepo    repositories.CartRepository
	productRepo repositories.ProductRepository
}

func NewCartService(cartRepo repositories.CartRepository, productRepo repositories.ProductRepository) CartService {
	return &cartService{cartRepo: cartRepo, productRepo: productRepo}
}

// Implement the CartService methods here (AddToCart, UpdateCartItem, RemoveFromCart, GetCart)
// Each method will interact with the cartRepo to perform the necessary operations and return the appropriate responses
func (s *cartService) AddToCart(userID uint, req *requests.AddToCartRequest) error {
	if _, err := s.productRepo.FindByID(req.ProductID); err != nil {
		return configs.ProductNotFound
	}

	existingItem, err := s.cartRepo.FindByUserAndProduct(userID, req.ProductID)
	if err == nil {
		existingItem.Quantity += req.Quantity
		if err := s.cartRepo.Update(existingItem); err != nil {
			return configs.UpdateCartItemFailed
		}
		return nil
	}

	newItem := &entities.CartItem{
		UserID:    userID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}

	if err := s.cartRepo.Create(newItem); err != nil {
		return configs.AddToCartFailed
	}

	return nil
}

// GetCart implements [CartService].
func (s *cartService) GetCart(userID uint) (*responses.CartResponse, error) {
	cartItems, err := s.cartRepo.ListByUserID(userID)
	if err != nil {
		return nil, err
	}

	var responseItems []responses.CartItemData
	var totalAmount float64

	if len(cartItems) == 0 {
		return &responses.CartResponse{
			Items:       []responses.CartItemData{},
			TotalAmount: 0,
		}, nil
	}

	for _, item := range cartItems {
		product, err := s.productRepo.FindByID(item.ProductID)
		if err != nil {
			continue // Skip items with missing products
		}

		subTotal := float64(item.Quantity) * product.Price
		totalAmount += subTotal

		responseItems = append(responseItems, responses.CartItemData{
			ID:        item.ID,
			ProductID: item.ProductID,
			Name:      product.Name,
			Price:     product.Price,
			Image:     product.Image,
			Quantity:  item.Quantity,
			SubTotal:  subTotal,
		})
	}

	return &responses.CartResponse{
		Items:       responseItems,
		TotalAmount: totalAmount,
	}, nil
}

// UpdateCartItem implements [CartService].
func (s *cartService) UpdateCartItem(userID uint, cartItemID uint, req *requests.UpdateCartItemRequest) error {
	cartItem, err := s.cartRepo.FindByID(cartItemID)
	if err != nil {
		return configs.CartItemNotFound
	}

	if cartItem.UserID != userID {
		return configs.AccessDenied
	}

	cartItem.Quantity = req.Quantity

	if err := s.cartRepo.Update(cartItem); err != nil {
		return configs.UpdateCartItemFailed
	}

	return nil
}

// RemoveFromCart implements [CartService].
func (s *cartService) RemoveFromCart(userID uint, cartItemID uint) error {
	cartItem, err := s.cartRepo.FindByID(cartItemID)
	if err != nil {
		return configs.CartItemNotFound
	}

	if cartItem.UserID != userID {
		return configs.AccessDenied
	}

	if err := s.cartRepo.Delete(cartItem); err != nil {
		return configs.RemoveFromCartFailed
	}

	return nil
}
