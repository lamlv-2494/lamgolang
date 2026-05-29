package services

import (
	"food_delivery/internal/configs"
	"food_delivery/internal/models/dto/requests"
	"food_delivery/internal/models/dto/responses"
	"food_delivery/internal/models/entities"
	"food_delivery/internal/repositories"
)

type CategoryService interface {
	CreateCategory(req requests.CreateCategoryRequest) (*responses.CategoryResponse, error)
	UpdateCategory(id uint, req requests.UpdateCategoryRequest) (*responses.CategoryResponse, error)
	DeleteCategory(id uint) error
	GetCategories(page, limit int) (responses.ListResponse[*responses.CategoryResponse], error)
}

type categoryService struct {
	repo repositories.CategoryRepository
}

func NewCategoryService(repo repositories.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) CreateCategory(req requests.CreateCategoryRequest) (*responses.CategoryResponse, error) {
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

	return &responses.CategoryResponse{
		ID:          int(category.ID),
		Name:        category.Name,
		Description: category.Description,
	}, nil
}

func (s *categoryService) UpdateCategory(id uint, req requests.UpdateCategoryRequest) (*responses.CategoryResponse, error) {
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

	return &responses.CategoryResponse{
		ID:          int(category.ID),
		Name:        category.Name,
		Description: category.Description,
	}, nil
}

func (s *categoryService) DeleteCategory(id uint) error {
	category, err := s.repo.FindByID(id)
	if err != nil {
		return configs.CategoryNotFound
	}

	return s.repo.Delete(category)
}

func (s *categoryService) GetCategories(page, limit int) (responses.ListResponse[*responses.CategoryResponse], error) {
	categories, totalCount, err := s.repo.List(page, limit)
	if err != nil {
		return responses.ListResponse[*responses.CategoryResponse]{}, err
	}

	categoryResponses := make([]*responses.CategoryResponse, len(categories))
	for i, category := range categories {
		categoryResponses[i] = &responses.CategoryResponse{
			ID:          int(category.ID),
			Name:        category.Name,
			Description: category.Description,
		}
	}

	return responses.ListResponse[*responses.CategoryResponse]{
		Items:      categoryResponses,
		TotalCount: totalCount,
	}, nil
}
