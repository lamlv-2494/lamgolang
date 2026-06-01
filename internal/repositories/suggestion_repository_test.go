package repositories

import (
	"testing"

	"food_delivery/internal/models/entities"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSuggestionRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSuggestionRepository(db)
	userRepo := NewUserRepository(db)

	user := &entities.User{Username: "suggester1", Email: "s1@example.com", Password: "pass"}
	require.NoError(t, userRepo.CreateUser(user))

	suggestion := &entities.Suggestion{
		UserID:      user.ID,
		Title:       "Add more vegan options",
		Description: "Please add more plant-based dishes",
	}

	err := repo.Create(suggestion)
	require.NoError(t, err)
	assert.NotZero(t, suggestion.ID)
}

func TestSuggestionRepository_ListAll(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSuggestionRepository(db)
	userRepo := NewUserRepository(db)

	user := &entities.User{Username: "suggester2", Email: "s2@example.com", Password: "pass"}
	require.NoError(t, userRepo.CreateUser(user))

	for i := 0; i < 5; i++ {
		s := &entities.Suggestion{
			UserID:      user.ID,
			Title:       "Suggestion " + string(rune('A'+i)),
			Description: "Description " + string(rune('A'+i)),
		}
		require.NoError(t, repo.Create(s))
	}

	suggestions, total, err := repo.ListAll(1, 3)
	require.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, suggestions, 3)

	suggestions2, total2, err := repo.ListAll(2, 3)
	require.NoError(t, err)
	assert.Equal(t, int64(5), total2)
	assert.Len(t, suggestions2, 2)
}

func TestSuggestionRepository_ListAll_Empty(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSuggestionRepository(db)

	suggestions, total, err := repo.ListAll(1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Empty(t, suggestions)
}
