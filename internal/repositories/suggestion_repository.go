package repositories

import (
	"food_delivery/internal/models/entities"

	"gorm.io/gorm"
)

type SuggestionRepository interface {
	Create(suggestion *entities.Suggestion) error
	ListAll() ([]*entities.Suggestion, error)
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

func (r *suggestionRepository) ListAll() ([]*entities.Suggestion, error) {
	var suggestions []*entities.Suggestion
	err := r.db.
		Preload("User").
		Order("created_at desc").
		Find(&suggestions).
		Error
	if err != nil {
		return nil, err
	}
	return suggestions, nil
}
