package services

import (
	"food_delivery/internal/configs"
	"food_delivery/internal/models/dto/requests"
	"food_delivery/internal/models/dto/responses"
	"food_delivery/internal/models/entities"
	"food_delivery/internal/repositories"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(req requests.RegisterRequest) (*responses.UserResponse, error)
	Login(req requests.LoginRequest) (*responses.UserResponse, error)
	GetCurrentUser(id uint) (*responses.UserResponse, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (u *userService) Register(req requests.RegisterRequest) (*responses.UserResponse, error) {
	if ex, _ := u.repo.FindByEmail(req.User.Email); ex != nil {
		return nil, configs.EmailOrUserTaken
	}
	if ex, _ := u.repo.FindByUsername(req.User.Username); ex != nil {
		return nil, configs.EmailOrUserTaken
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.User.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, configs.CreateUserFailed
	}

	user := &entities.User{
		Username: req.User.Username,
		Email:    req.User.Email,
		Password: string(hashedPassword),
		Role:     "user",
	}

	if err := u.repo.CreateUser(user); err != nil {
		return nil, configs.CreateUserFailed
	}

	token, err := configs.GenerateToken(user.ID, user.Role)

	if err != nil {
		return nil, configs.GenerateTokenFail
	}

	userResponse := &responses.UserResponse{
		User: responses.UserData{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
			Token:    token,
			Role:     user.Role,
		},
	}

	return userResponse, nil
}

func (u *userService) Login(req requests.LoginRequest) (*responses.UserResponse, error) {
	user, err := u.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, configs.UserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, configs.InvalidPassword
	}

	token, err := configs.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, configs.GenerateTokenFail
	}

	userResponse := &responses.UserResponse{
		User: responses.UserData{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
			Token:    token,
			Role:     user.Role,
		},
	}
	return userResponse, nil
}

func (u *userService) GetCurrentUser(id uint) (*responses.UserResponse, error) {
	user, err := u.repo.FindByID(id)
	if err != nil {
		return nil, configs.UserNotFound
	}

	userResponse := &responses.UserResponse{
		User: responses.UserData{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
			Role:     user.Role,
		},
	}
	return userResponse, nil
}
