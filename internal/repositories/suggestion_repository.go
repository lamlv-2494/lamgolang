package repositories

import (
	"food_delivery/internal/models/entities"

	"gorm.io/gorm"
)

type SuggestionRepository interface {
	Create(suggestion *entities.Suggestion) error
	ListAll(page, limit int) ([]*entities.Suggestion, int64, error)
}

type suggestionRepository struct {
	db *gorm.DB
}

func NewSuggestionRepository(db *gorm.DB) SuggestionRepository {
	return &suggestionRepository{db: db}
}

func (r *suggestionRepository) Create(suggestion *entities.Suggestion) error {
	return r.db.Create(suggestion).Error
}

func (r *suggestionRepository) ListAll(page, limit int) ([]*entities.Suggestion, int64, error) {
	var suggestions []*entities.Suggestion
	var totalCount int64

	err := r.db.Model(&entities.Suggestion{}).Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.
		Preload("User").
		Order("created_at desc").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&suggestions).
		Error
	if err != nil {
		return nil, totalCount, err
	}
	return suggestions, totalCount, nil
}
