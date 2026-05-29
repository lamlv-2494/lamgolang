package services

import (
	"food_delivery/internal/configs"
	"food_delivery/internal/models/dto/requests"
	"food_delivery/internal/models/dto/responses"
	"food_delivery/internal/models/entities"
	"food_delivery/internal/repositories"
)

type RatingService interface {
	CreateRating(userID, productID uint, req requests.CreateRatingRequest) error
	GetRatingByID(userID uint, productID uint) (*responses.RatingResponse, error)

	GetRatingsByUserID(userID uint, page, limit int) (*responses.ListResponse[*responses.RatingResponseData], error)
}

type ratingService struct {
	ratingRepo  repositories.RatingRepository
	productRepo repositories.ProductRepository
}

func NewRatingService(ratingRepo repositories.RatingRepository, productRepo repositories.ProductRepository) RatingService {
	return &ratingService{ratingRepo: ratingRepo, productRepo: productRepo}
}

func (s *ratingService) CreateRating(userID, productID uint, req requests.CreateRatingRequest) error {
	if _, err := s.productRepo.FindByID(productID); err != nil {
		return configs.ProductNotFound
	}

	if _, err := s.ratingRepo.FindByUserAndProduct(userID, productID); err == nil {
		return configs.RatingAlreadyExists
	}

	rating := &entities.Rating{
		UserID:    userID,
		ProductID: productID,
		Stars:     req.Stars,
		Comment:   req.Comment,
	}

	if err := s.ratingRepo.Create(rating); err != nil {
		return configs.CreateRatingFailed
	}

	return nil
}

func (s *ratingService) GetRatingByID(userID uint, productID uint) (*responses.RatingResponse, error) {
	rating, err := s.ratingRepo.FindByUserAndProduct(userID, productID)
	if err != nil {
		return nil, configs.RatingNotFound
	}

	if rating == nil {
		return &responses.RatingResponse{}, nil
	}

	return &responses.RatingResponse{
		Rating: responses.RatingResponseData{
			ProductID: rating.ProductID,
			Stars:     rating.Stars,
			Comment:   rating.Comment,
		},
	}, nil
}

func (s *ratingService) GetRatingsByUserID(userID uint, page, limit int) (*responses.ListResponse[*responses.RatingResponseData], error) {

	ratings, totalCount, err := s.ratingRepo.GetRatingsByUserID(userID, page, limit)
	if err != nil {
		return nil, configs.RatingNotFound
	}

	if ratings == nil {
		return &responses.ListResponse[*responses.RatingResponseData]{
			Items:      []*responses.RatingResponseData{},
			TotalCount: totalCount,
		}, nil
	}

	var ratingResponses []*responses.RatingResponseData
	for _, rating := range ratings {
		ratingResponses = append(ratingResponses, &responses.RatingResponseData{
			ProductID: rating.ProductID,
			Stars:     rating.Stars,
			Comment:   rating.Comment,
		})
	}

	return &responses.ListResponse[*responses.RatingResponseData]{
		Items:      ratingResponses,
		TotalCount: totalCount,
	}, nil
}
