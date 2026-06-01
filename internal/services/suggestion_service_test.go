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

func newSuggestionRepo() *mockSuggestionRepo {
	return &mockSuggestionRepo{
		createFn: func(s *entities.Suggestion) error { return nil },
		listAllFn: func(p, l int) ([]*entities.Suggestion, int64, error) {
			return nil, 0, nil
		},
	}
}

// ─── CreateSuggestion ─────────────────────────────────────────────────────────

func TestSuggestionService_Create_Success(t *testing.T) {
	repo := newSuggestionRepo()
	svc := NewSuggestionService(repo)

	err := svc.CreateSuggestion(1, requests.CreateSuggestionRequest{
		Title:       "Better Menu",
		Description: "Add more options",
	})
	assert.NoError(t, err)
}

func TestSuggestionService_Create_Fails(t *testing.T) {
	repo := newSuggestionRepo()
	repo.createFn = func(s *entities.Suggestion) error { return errors.New("db error") }
	svc := NewSuggestionService(repo)

	err := svc.CreateSuggestion(1, requests.CreateSuggestionRequest{Title: "Test"})
	assert.Equal(t, configs.CreateSuggestionFailed, err)
}

// ─── GetAllSuggestions ────────────────────────────────────────────────────────

func TestSuggestionService_GetAll_Success(t *testing.T) {
	repo := newSuggestionRepo()
	repo.listAllFn = func(p, l int) ([]*entities.Suggestion, int64, error) {
		s := &entities.Suggestion{
			Title:       "More Drinks",
			Description: "Please add juice",
			User:        entities.User{Username: "alice", Email: "alice@example.com"},
		}
		s.UserID = 1
		return []*entities.Suggestion{s}, 1, nil
	}
	svc := NewSuggestionService(repo)

	resp, err := svc.GetAllSuggestions(1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), resp.TotalCount)
	assert.Len(t, resp.Items, 1)
	assert.Equal(t, "More Drinks", resp.Items[0].Title)
	assert.Equal(t, "alice", resp.Items[0].Username)
}

func TestSuggestionService_GetAll_Empty(t *testing.T) {
	repo := newSuggestionRepo()
	svc := NewSuggestionService(repo)

	resp, err := svc.GetAllSuggestions(1, 10)
	require.NoError(t, err)
	assert.Empty(t, resp.Items)
	assert.Equal(t, int64(0), resp.TotalCount)
}

func TestSuggestionService_GetAll_Fails(t *testing.T) {
	repo := newSuggestionRepo()
	repo.listAllFn = func(p, l int) ([]*entities.Suggestion, int64, error) {
		return nil, 0, errors.New("db error")
	}
	svc := NewSuggestionService(repo)

	_, err := svc.GetAllSuggestions(1, 10)
	assert.Equal(t, configs.FetchSuggestionsFailed, err)
}
