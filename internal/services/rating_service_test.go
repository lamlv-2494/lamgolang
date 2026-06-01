package services

import (
	"errors"
	"testing"

	"food_delivery/internal/configs"
	"food_delivery/internal/models/dto/requests"
	"food_delivery/internal/models/entities"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newRatingRepo() *mockRatingRepo {
	return &mockRatingRepo{
		createFn: func(r *entities.Rating) error { return nil },
		findByUserAndProductFn: func(uID, pID uint) (*entities.Rating, error) {
			return nil, errors.New("not found")
		},
		getRatingsByUserIDFn: func(uID uint, p, l int) ([]*entities.Rating, int64, error) {
			return nil, 0, nil
		},
	}
}

// ─── CreateRating ─────────────────────────────────────────────────────────────

func TestRatingService_CreateRating_Success(t *testing.T) {
	prodRepo := newProdRepo()
	prodRepo.findByIDFn = func(id uint) (*entities.Product, error) { return baseProduct(id), nil }
	ratingRepo := newRatingRepo()
	svc := NewRatingService(ratingRepo, prodRepo)

	err := svc.CreateRating(1, 1, requests.CreateRatingRequest{Stars: 5, Comment: "Great!"})
	assert.NoError(t, err)
}

func TestRatingService_CreateRating_ProductNotFound(t *testing.T) {
	prodRepo := newProdRepo()
	ratingRepo := newRatingRepo()
	svc := NewRatingService(ratingRepo, prodRepo)

	err := svc.CreateRating(1, 99, requests.CreateRatingRequest{Stars: 4})
	assert.Equal(t, configs.ProductNotFound, err)
}

func TestRatingService_CreateRating_AlreadyRated(t *testing.T) {
	prodRepo := newProdRepo()
	prodRepo.findByIDFn = func(id uint) (*entities.Product, error) { return baseProduct(id), nil }
	ratingRepo := newRatingRepo()
	ratingRepo.findByUserAndProductFn = func(uID, pID uint) (*entities.Rating, error) {
		return &entities.Rating{Stars: 3}, nil // rating already exists
	}
	svc := NewRatingService(ratingRepo, prodRepo)

	err := svc.CreateRating(1, 1, requests.CreateRatingRequest{Stars: 5})
	assert.Equal(t, configs.RatingAlreadyExists, err)
}

func TestRatingService_CreateRating_CreateFails(t *testing.T) {
	prodRepo := newProdRepo()
	prodRepo.findByIDFn = func(id uint) (*entities.Product, error) { return baseProduct(id), nil }
	ratingRepo := newRatingRepo()
	ratingRepo.createFn = func(r *entities.Rating) error { return errors.New("db error") }
	svc := NewRatingService(ratingRepo, prodRepo)

	err := svc.CreateRating(1, 1, requests.CreateRatingRequest{Stars: 4})
	assert.Equal(t, configs.CreateRatingFailed, err)
}

// ─── GetRatingByID ────────────────────────────────────────────────────────────

func TestRatingService_GetRatingByID_Found(t *testing.T) {
	prodRepo := newProdRepo()
	ratingRepo := newRatingRepo()
	ratingRepo.findByUserAndProductFn = func(uID, pID uint) (*entities.Rating, error) {
		return &entities.Rating{ProductID: 1, Stars: 4, Comment: "Nice"}, nil
	}
	svc := NewRatingService(ratingRepo, prodRepo)

	resp, err := svc.GetRatingByID(1, 1)
	require.NoError(t, err)
	assert.Equal(t, 4, resp.Rating.Stars)
	assert.Equal(t, "Nice", resp.Rating.Comment)
}

func TestRatingService_GetRatingByID_NotFound(t *testing.T) {
	prodRepo := newProdRepo()
	ratingRepo := newRatingRepo()
	svc := NewRatingService(ratingRepo, prodRepo)

	_, err := svc.GetRatingByID(1, 99)
	assert.Equal(t, configs.RatingNotFound, err)
}

func TestRatingService_GetRatingByID_NilRating(t *testing.T) {
	prodRepo := newProdRepo()
	ratingRepo := newRatingRepo()
	ratingRepo.findByUserAndProductFn = func(uID, pID uint) (*entities.Rating, error) {
		return nil, nil // returns nil rating with no error
	}
	svc := NewRatingService(ratingRepo, prodRepo)

	resp, err := svc.GetRatingByID(1, 1)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 0, resp.Rating.Stars) // empty response
}

// ─── GetRatingsByUserID ───────────────────────────────────────────────────────

func TestRatingService_GetRatingsByUserID_Success(t *testing.T) {
	prodRepo := newProdRepo()
	ratingRepo := newRatingRepo()
	prod := baseProduct(1)
	ratingRepo.getRatingsByUserIDFn = func(uID uint, p, l int) ([]*entities.Rating, int64, error) {
		return []*entities.Rating{
			{ProductID: 1, Stars: 5, Comment: "Excellent", Product: prod},
		}, 1, nil
	}
	svc := NewRatingService(ratingRepo, prodRepo)

	resp, err := svc.GetRatingsByUserID(1, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), resp.TotalCount)
	assert.Len(t, resp.Items, 1)
	assert.Equal(t, "Burger", resp.Items[0].ProductName)
}

func TestRatingService_GetRatingsByUserID_NilProductInRating(t *testing.T) {
	prodRepo := newProdRepo()
	ratingRepo := newRatingRepo()
	ratingRepo.getRatingsByUserIDFn = func(uID uint, p, l int) ([]*entities.Rating, int64, error) {
		return []*entities.Rating{
			{ProductID: 1, Stars: 3, Product: nil},
		}, 1, nil
	}
	svc := NewRatingService(ratingRepo, prodRepo)

	resp, err := svc.GetRatingsByUserID(1, 1, 10)
	require.NoError(t, err)
	assert.Len(t, resp.Items, 1)
	assert.Equal(t, "", resp.Items[0].ProductName) // empty when product is nil
}

func TestRatingService_GetRatingsByUserID_NilSlice(t *testing.T) {
	prodRepo := newProdRepo()
	ratingRepo := newRatingRepo()
	ratingRepo.getRatingsByUserIDFn = func(uID uint, p, l int) ([]*entities.Rating, int64, error) {
		return nil, 0, nil
	}
	svc := NewRatingService(ratingRepo, prodRepo)

	resp, err := svc.GetRatingsByUserID(1, 1, 10)
	require.NoError(t, err)
	assert.Empty(t, resp.Items)
}

func TestRatingService_GetRatingsByUserID_Fails(t *testing.T) {
	prodRepo := newProdRepo()
	ratingRepo := newRatingRepo()
	ratingRepo.getRatingsByUserIDFn = func(uID uint, p, l int) ([]*entities.Rating, int64, error) {
		return nil, 0, errors.New("db error")
	}
	svc := NewRatingService(ratingRepo, prodRepo)

	_, err := svc.GetRatingsByUserID(1, 1, 10)
	assert.Equal(t, configs.RatingNotFound, err)
}
