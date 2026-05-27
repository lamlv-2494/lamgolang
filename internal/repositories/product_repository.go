package repositories

import (
	"food_delivery/internal/models/entities"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product *entities.Product) error
	FindByID(id uint) (*entities.Product, error)
	Update(product *entities.Product) error
	Delete(product *entities.Product) error
	List() ([]*entities.Product, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *entities.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) FindByID(id uint) (*entities.Product, error) {
	var product entities.Product
	if err := r.db.Preload("Category").Preload("Ratings").First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) Update(product *entities.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(product *entities.Product) error {
	return r.db.Delete(product).Error
}

func (r *productRepository) List() ([]*entities.Product, error) {
	var products []*entities.Product
	if err := r.db.Preload("Category").Preload("Ratings").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}
