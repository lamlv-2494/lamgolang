package repositories

import (
	"food_delivery/internal/models/entities"
	"strings"

	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product *entities.Product) error
	FindByID(id uint) (*entities.Product, error)
	Update(product *entities.Product) error
	Delete(product *entities.Product) error
	List(classify string, categoryID uint, minPrice, maxPrice float64, minRating float64, sort string) ([]*entities.Product, error)
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

func (r *productRepository) List(classify string, categoryID uint, minPrice, maxPrice float64, minRating float64, sort string) ([]*entities.Product, error) {
	var products []*entities.Product

	query := r.db.Preload("Category").Preload("Ratings").Model(&entities.Product{})

	if classify != "" {
		query = query.Where("category_id IN (SELECT id FROM categories WHERE LOWER(name) LIKE ?)", "%"+strings.ToLower(classify)+"%")
	}

	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}

	if minPrice > 0 {
		query = query.Where("price >= ?", minPrice)
	}

	if maxPrice > 0 {
		query = query.Where("price <= ?", maxPrice)
	}

	if minRating > 0 {
		query = query.Having("AVG(ratings.value) >= ?", minRating)
	}

	switch sort {
	case "name_asc":
		query = query.Order("name ASC") // Alphabet A-Z
	case "name_desc":
		query = query.Order("name DESC") // Alphabet Z-A
	case "price_asc":
		query = query.Order("price ASC") // Price low to high
	case "price_desc":
		query = query.Order("price DESC") // Price high to low
	default:
		query = query.Order("id DESC") // Default sort
	}

	if err := query.Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}
