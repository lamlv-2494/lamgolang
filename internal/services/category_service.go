package services

import (
	"food_delivery/internal/configs"
	"food_delivery/internal/models/dto/requests"
	"food_delivery/internal/models/dto/responses"
	"food_delivery/internal/models/entities"
	"food_delivery/internal/repositories"
)

type CategoryService interface {
	CreateCategory(req requests.CreateCategoryRequest) (*entities.Category, error)
	UpdateCategory(id uint, req requests.UpdateCategoryRequest) (*entities.Category, error)
	DeleteCategory(id uint) error
	GetCategories(page, limit int) (responses.AnyListResponse, error)
}

type categoryService struct {
	repo repositories.CategoryRepository
}

func NewCategoryService(repo repositories.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) CreateCategory(req requests.CreateCategoryRequest) (*entities.Category, error) {
	if ex, _ := s.repo.FindByName(req.Name); ex != nil {
		return nil, configs.CategoryAlreadyExists
	}

	category := &entities.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.repo.Create(category); err != nil {
		return nil, configs.CreateCategoryFailed
	}

	return category, nil
}

func (s *categoryService) UpdateCategory(id uint, req requests.UpdateCategoryRequest) (*entities.Category, error) {
	category, err := s.repo.FindByID(id)
	if err != nil {
		return nil, configs.CategoryNotFound
	}

	if req.Name != nil {
		category.Name = *req.Name
	}
	if req.Description != nil {
		category.Description = *req.Description
	}

	if err := s.repo.Update(category); err != nil {
		return nil, configs.UpdateCategoryFailed
	}

	return category, nil
}

func (s *categoryService) DeleteCategory(id uint) error {
	category, err := s.repo.FindByID(id)
	if err != nil {
		return configs.CategoryNotFound
	}

	return s.repo.Delete(category)
}

func (s *categoryService) GetCategories(page, limit int) (responses.AnyListResponse, error) {
	categories, totalCount, err := s.repo.List(page, limit)
	if err != nil {
		return responses.AnyListResponse{}, err
	}

	return responses.AnyListResponse{
		Data:       categories,
		TotalCount: totalCount,
	}, nil
}
