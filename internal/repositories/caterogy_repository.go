package repositories

import (
	"food_delivery/internal/models/entities"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(category *entities.Category) error
	FindByID(id uint) (*entities.Category, error)
	FindByName(name string) (*entities.Category, error)
	Update(category *entities.Category) error
	Delete(category *entities.Category) error
	List() ([]*entities.Category, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(category *entities.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) FindByID(id uint) (*entities.Category, error) {
	var category entities.Category
	if err := r.db.First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) FindByName(name string) (*entities.Category, error) {
	var category entities.Category
	if err := r.db.Where("name = ?", name).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) Update(category *entities.Category) error {
	return r.db.Save(category).Error
}

func (r *categoryRepository) Delete(category *entities.Category) error {
	return r.db.Delete(category).Error
}

func (r *categoryRepository) List() ([]*entities.Category, error) {
	var categories []*entities.Category
	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}
