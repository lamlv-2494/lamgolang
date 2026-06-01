package repositories

import (
	"testing"

	"food_delivery/internal/models/entities"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCategoryRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewCategoryRepository(db)

	cat := &entities.Category{Name: "Burgers", Description: "Tasty burgers"}
	err := repo.Create(cat)
	require.NoError(t, err)
	assert.NotZero(t, cat.ID)
}

func TestCategoryRepository_FindByID_Found(t *testing.T) {
	db := setupTestDB(t)
	repo := NewCategoryRepository(db)

	cat := &entities.Category{Name: "Pizza", Description: "Italian pizza"}
	require.NoError(t, repo.Create(cat))

	found, err := repo.FindByID(cat.ID)
	require.NoError(t, err)
	assert.Equal(t, "Pizza", found.Name)
}

func TestCategoryRepository_FindByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewCategoryRepository(db)

	_, err := repo.FindByID(9999)
	assert.Error(t, err)
}

func TestCategoryRepository_FindByName_Found(t *testing.T) {
	db := setupTestDB(t)
	repo := NewCategoryRepository(db)

	cat := &entities.Category{Name: "Sushi", Description: "Japanese sushi"}
	require.NoError(t, repo.Create(cat))

	found, err := repo.FindByName("Sushi")
	require.NoError(t, err)
	assert.Equal(t, cat.ID, found.ID)
}

func TestCategoryRepository_FindByName_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewCategoryRepository(db)

	_, err := repo.FindByName("Unknown")
	assert.Error(t, err)
}

func TestCategoryRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewCategoryRepository(db)

	cat := &entities.Category{Name: "Tacos", Description: "Mexican tacos"}
	require.NoError(t, repo.Create(cat))

	cat.Description = "Authentic Mexican tacos"
	err := repo.Update(cat)
	require.NoError(t, err)

	found, err := repo.FindByID(cat.ID)
	require.NoError(t, err)
	assert.Equal(t, "Authentic Mexican tacos", found.Description)
}

func TestCategoryRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewCategoryRepository(db)

	cat := &entities.Category{Name: "ToDelete", Description: "will be deleted"}
	require.NoError(t, repo.Create(cat))

	err := repo.Delete(cat)
	require.NoError(t, err)

	_, err = repo.FindByID(cat.ID)
	assert.Error(t, err)
}

func TestCategoryRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := NewCategoryRepository(db)

	for i := 0; i < 5; i++ {
		c := &entities.Category{Name: "Cat" + string(rune('A'+i))}
		require.NoError(t, repo.Create(c))
	}

	cats, total, err := repo.List(1, 3)
	require.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, cats, 3)

	cats2, total2, err := repo.List(2, 3)
	require.NoError(t, err)
	assert.Equal(t, int64(5), total2)
	assert.Len(t, cats2, 2)
}

func TestCategoryRepository_List_Empty(t *testing.T) {
	db := setupTestDB(t)
	repo := NewCategoryRepository(db)

	cats, total, err := repo.List(1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Empty(t, cats)
}
