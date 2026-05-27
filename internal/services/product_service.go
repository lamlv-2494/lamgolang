package services

import (
	"food_delivery/internal/configs"
	"food_delivery/internal/models/dto/requests"
	"food_delivery/internal/models/dto/responses"
	"food_delivery/internal/models/entities"
	"food_delivery/internal/repositories"
)

type ProductService interface {
	CreateProduct(req requests.CreateProductRequest) (*responses.ProductData, error)
	UpdateProduct(id uint, req requests.UpdateProductRequest) (*responses.ProductData, error)
	DeleteProduct(id uint) error
	GetProducts(classify string, categoryID uint, minPrice, maxPrice float64, minRating float64, sort string) ([]*responses.ProductData, error)

	GetProductByID(id uint) (*responses.ProductData, error)
}

type productService struct {
	productRepo  repositories.ProductRepository
	categoryRepo repositories.CategoryRepository
}

func NewProductService(productRepo repositories.ProductRepository, categoryRepo repositories.CategoryRepository) ProductService {
	return &productService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *productService) CreateProduct(req requests.CreateProductRequest) (*responses.ProductData, error) {
	// Kiểm tra xem danh mục món ăn truyền lên có tồn tại không
	if _, err := s.categoryRepo.FindByID(req.CategoryID); err != nil {
		return nil, configs.CategoryNotFound
	}

	product := &entities.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Image:       req.Image,
		CategoryID:  req.CategoryID,
	}

	if err := s.productRepo.Create(product); err != nil {
		return nil, configs.CreateProductFailed
	}

	newProduct, err := s.productRepo.FindByID(product.ID)
	if err != nil {
		return nil, configs.ProductNotFound
	}

	response := &responses.ProductData{
		ID:          newProduct.ID,
		Name:        newProduct.Name,
		Description: newProduct.Description,
		Price:       newProduct.Price,
		Image:       newProduct.Image,
		CategoryID:  newProduct.CategoryID,
		Category: responses.CategoryCompact{
			Name:        newProduct.Category.Name,
			Description: newProduct.Category.Description,
		},
	}

	return response, nil
}

func (s *productService) UpdateProduct(id uint, req requests.UpdateProductRequest) (*responses.ProductData, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, configs.ProductNotFound
	}

	if req.CategoryID != nil {
		if _, err := s.categoryRepo.FindByID(*req.CategoryID); err != nil {
			return nil, configs.CategoryNotFound
		}
		product.CategoryID = *req.CategoryID
	}

	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.Image != nil {
		product.Image = *req.Image
	}

	if err := s.productRepo.Update(product); err != nil {
		return nil, configs.UpdateProductFailed
	}

	newProduct, err := s.productRepo.FindByID(product.ID)
	if err != nil {
		return nil, configs.ProductNotFound
	}

	response := &responses.ProductData{
		ID:          newProduct.ID,
		Name:        newProduct.Name,
		Description: newProduct.Description,
		Price:       newProduct.Price,
		Image:       newProduct.Image,
		CategoryID:  newProduct.CategoryID,
		Category: responses.CategoryCompact{
			Name:        newProduct.Category.Name,
			Description: newProduct.Category.Description,
		},
	}

	return response, nil
}

func (s *productService) DeleteProduct(id uint) error {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return configs.ProductNotFound
	}

	if err := s.productRepo.Delete(product); err != nil {
		return configs.DeleteProductFailed
	}

	return nil
}

func (s *productService) GetProducts(classify string, categoryID uint, minPrice, maxPrice float64, minRating float64, sort string) ([]*responses.ProductData, error) {
	products, err := s.productRepo.List(classify, categoryID, minPrice, maxPrice, minRating, sort)
	if err != nil {
		return nil, configs.ProductNotFound
	}

	if len(products) == 0 {
		return []*responses.ProductData{}, nil
	}

	var productResponses []*responses.ProductData
	for _, p := range products {
		var totalStars int
		var avgRating float64
		if len(p.Ratings) > 0 {
			for _, r := range p.Ratings {
				totalStars += r.Stars
			}
			avgRating = float64(totalStars) / float64(len(p.Ratings))
		}

		productResponses = append(productResponses, &responses.ProductData{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Image:       p.Image,
			CategoryID:  p.CategoryID,
			Category: responses.CategoryCompact{
				Name:        p.Category.Name,
				Description: p.Category.Description,
			},
			Rating: avgRating,
		})
	}

	return productResponses, nil
}

func (s *productService) GetProductByID(id uint) (*responses.ProductData, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, configs.ProductNotFound
	}

	var totalStars int
	var avgRating float64
	if len(product.Ratings) > 0 {
		for _, r := range product.Ratings {
			totalStars += r.Stars
		}
		avgRating = float64(totalStars) / float64(len(product.Ratings))
	}

	productResponse := &responses.ProductData{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Image:       product.Image,
		CategoryID:  product.CategoryID,
		Category: responses.CategoryCompact{
			Name:        product.Category.Name,
			Description: product.Category.Description,
		},
		Rating: avgRating,
	}

	return productResponse, nil
}
