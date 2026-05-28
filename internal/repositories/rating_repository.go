package repositories

import (
	"food_delivery/internal/models/entities"

	"gorm.io/gorm"
)

type RatingRepository interface {
	Create(rating *entities.Rating) error
	FindByUserAndProduct(userID, productID uint) (*entities.Rating, error)

	GetRatingsByUserID(userID uint) ([]*entities.Rating, error)
}

type ratingRepository struct {
	db *gorm.DB
}

func NewRatingRepository(db *gorm.DB) RatingRepository {
	return &ratingRepository{db: db}
}

func (r *ratingRepository) Create(rating *entities.Rating) error {
	return r.db.Create(rating).Error
}

func (r *ratingRepository) FindByUserAndProduct(userID, productID uint) (*entities.Rating, error) {
	var rating entities.Rating
	err := r.db.
		Where("user_id = ? AND product_id = ?", userID, productID).
		First(&rating).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &rating, nil
}

func (r *ratingRepository) GetRatingsByUserID(userID uint) ([]*entities.Rating, error) {
	var ratings []*entities.Rating
	err := r.db.Where("user_id = ?", userID).Find(&ratings).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return ratings, nil
}
