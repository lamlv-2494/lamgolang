package repositories

import (
	"food_delivery/internal/models/entities"

	"gorm.io/gorm"
)

type RatingRepository interface {
	Create(rating *entities.Rating) error
	FindByUserAndProduct(userID, productID uint) (*entities.Rating, error)

	GetRatingsByUserID(userID uint, page, limit int) ([]*entities.Rating, int64, error)
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

func (r *ratingRepository) GetRatingsByUserID(userID uint, page, limit int) ([]*entities.Rating, int64, error) {
	var ratings []*entities.Rating
	var totalCount int64

	err := r.db.Model(&entities.Rating{}).Where("user_id = ?", userID).Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.
		Preload("Product").
		Where("user_id = ?", userID).
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&ratings).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, totalCount, nil
		}
		return nil, totalCount, err
	}
	return ratings, totalCount, nil
}
