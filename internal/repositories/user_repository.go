package repositories

import (
	"food_delivery/internal/models/entities"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *entities.User) error
	FindByEmail(email string) (*entities.User, error)
	FindByID(id uint) (*entities.User, error)
	FindByUsername(username string) (*entities.User, error)
	UpdateUser(user *entities.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (u *userRepository) CreateUser(user *entities.User) error {
	return u.db.Create(user).Error
}

func (u *userRepository) FindByEmail(email string) (*entities.User, error) {
	var user entities.User
	if err := u.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepository) FindByID(id uint) (*entities.User, error) {
	var user entities.User
	err := u.db.First(&user, id).Error
	return &user, err
}

func (u *userRepository) FindByUsername(username string) (*entities.User, error) {
	var user entities.User
	if err := u.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepository) UpdateUser(user *entities.User) error {
	return u.db.Save(user).Error
}
