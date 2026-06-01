package repositories

import (
	"testing"

	"food_delivery/internal/models/entities"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_CreateUser_Success(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &entities.User{
		Username: "john",
		Email:    "john@example.com",
		Password: "hashed_password",
	}

	err := repo.CreateUser(user)
	require.NoError(t, err)
	assert.NotZero(t, user.ID)
}

func TestUserRepository_CreateUser_DuplicateEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user1 := &entities.User{Username: "alice", Email: "dup@example.com", Password: "pass"}
	require.NoError(t, repo.CreateUser(user1))

	user2 := &entities.User{Username: "bob", Email: "dup@example.com", Password: "pass"}
	err := repo.CreateUser(user2)
	assert.Error(t, err)
}

func TestUserRepository_FindByEmail_Found(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &entities.User{Username: "alice", Email: "alice@example.com", Password: "pass"}
	require.NoError(t, repo.CreateUser(user))

	found, err := repo.FindByEmail("alice@example.com")
	require.NoError(t, err)
	assert.Equal(t, "alice", found.Username)
	assert.Equal(t, "alice@example.com", found.Email)
}

func TestUserRepository_FindByEmail_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	_, err := repo.FindByEmail("nobody@example.com")
	assert.Error(t, err)
}

func TestUserRepository_FindByID_Found(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &entities.User{Username: "bob", Email: "bob@example.com", Password: "pass"}
	require.NoError(t, repo.CreateUser(user))

	found, err := repo.FindByID(user.ID)
	require.NoError(t, err)
	assert.Equal(t, user.ID, found.ID)
	assert.Equal(t, "bob", found.Username)
}

func TestUserRepository_FindByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	_, err := repo.FindByID(9999)
	assert.Error(t, err)
}

func TestUserRepository_FindByUsername_Found(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &entities.User{Username: "carol", Email: "carol@example.com", Password: "pass"}
	require.NoError(t, repo.CreateUser(user))

	found, err := repo.FindByUsername("carol")
	require.NoError(t, err)
	assert.Equal(t, "carol@example.com", found.Email)
}

func TestUserRepository_FindByUsername_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	_, err := repo.FindByUsername("ghost")
	assert.Error(t, err)
}

func TestUserRepository_UpdateUser(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &entities.User{Username: "dave", Email: "dave@example.com", Password: "pass"}
	require.NoError(t, repo.CreateUser(user))

	user.Username = "dave_updated"
	err := repo.UpdateUser(user)
	require.NoError(t, err)

	found, err := repo.FindByID(user.ID)
	require.NoError(t, err)
	assert.Equal(t, "dave_updated", found.Username)
}

func TestUserRepository_FindAllUsers(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	for i := 0; i < 5; i++ {
		u := &entities.User{
			Username: "user" + string(rune('0'+i)),
			Email:    "user" + string(rune('0'+i)) + "@example.com",
			Password: "pass",
		}
		require.NoError(t, repo.CreateUser(u))
	}

	users, total, err := repo.FindAllUsers(1, 3)
	require.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, users, 3)

	// page 2
	users2, _, err := repo.FindAllUsers(2, 3)
	require.NoError(t, err)
	assert.Len(t, users2, 2)
}

func TestUserRepository_DeleteUser(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &entities.User{Username: "eve", Email: "eve@example.com", Password: "pass"}
	require.NoError(t, repo.CreateUser(user))

	err := repo.DeleteUser(user)
	require.NoError(t, err)

	_, err = repo.FindByID(user.ID)
	assert.Error(t, err)
}
