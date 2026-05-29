package services

import (
	"food_delivery/internal/configs"
	"food_delivery/internal/models/dto/requests"
	"food_delivery/internal/models/dto/responses"
	"food_delivery/internal/models/entities"
	"food_delivery/internal/repositories"
)

type SuggestionService interface {
	CreateSuggestion(userID uint, req requests.CreateSuggestionRequest) error
	GetAllSuggestions() (*responses.SuggestionListResponse, error)
}

type suggestionService struct {
	suggestionRepo repositories.SuggestionRepository
}

func NewSuggestionService(suggestionRepo repositories.SuggestionRepository) SuggestionService {
	return &suggestionService{suggestionRepo: suggestionRepo}
}

func (s *suggestionService) CreateSuggestion(userID uint, req requests.CreateSuggestionRequest) error {
	suggestion := &entities.Suggestion{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
	}

	if err := s.suggestionRepo.Create(suggestion); err != nil {
		return configs.CreateSuggestionFailed
	}
	return nil
}

func (s *suggestionService) GetAllSuggestions() (*responses.SuggestionListResponse, error) {
	suggestions, err := s.suggestionRepo.ListAll()
	if err != nil {
		return nil, configs.FetchSuggestionsFailed
	}

	var suggestionResponses []responses.SuggestionData
	for _, s := range suggestions {
		suggestionResponses = append(suggestionResponses, responses.SuggestionData{
			ID:          s.ID,
			UserID:      s.UserID,
			Username:    s.User.Username,
			Email:       s.User.Email,
			Title:       s.Title,
			Description: s.Description,
			CreatedAt:   s.CreatedAt,
		})
	}

	return &responses.SuggestionListResponse{Suggestions: suggestionResponses}, nil
}
