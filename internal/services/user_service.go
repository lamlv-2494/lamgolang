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
	UpdateUser(id uint, req *requests.UpdateUserRequest) (*responses.UserResponse, error)

	// Admin functions
	GetAllUsers(page, limit int) (*responses.AnyListResponse, error)
	AdminUpdateUser(targetID uint, req *requests.AdminUpdateUserRequest) (*responses.UserResponse, error)
	AdminDeleteUser(targetID uint) error
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

func (u *userService) UpdateUser(id uint, req *requests.UpdateUserRequest) (*responses.UserResponse, error) {
	user, err := u.repo.FindByID(id)
	if err != nil {
		return nil, configs.UserNotFound
	}

	if req.User.Username != nil {
		existingUser, err := u.repo.FindByUsername(*req.User.Username)
		if err == nil && existingUser.ID != id {
			return nil, configs.UsernameTaken
		}
		user.Username = *req.User.Username
	}
	if req.User.Email != nil {
		existingUser, err := u.repo.FindByEmail(*req.User.Email)
		if err == nil && existingUser.ID != id {
			return nil, configs.EmailTaken
		}
		user.Email = *req.User.Email
	}
	if req.User.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.User.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, configs.UpdateUserFailed
		}
		if user.Password != string(hashedPassword) {
			user.Password = string(hashedPassword)
		}
	}

	if err := u.repo.UpdateUser(user); err != nil {
		return nil, configs.UpdateUserFailed
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

// Admin functions
func (u *userService) GetAllUsers(page, limit int) (*responses.AnyListResponse, error) {
	users, totalCount, err := u.repo.FindAllUsers(page, limit)
	if err != nil {
		return nil, configs.FetchUsersFailed
	}

	var userResponses []responses.UserData
	for _, user := range users {
		userResponses = append(userResponses, responses.UserData{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
			Role:     user.Role,
		})
	}
	return &responses.AnyListResponse{
		Data:       userResponses,
		TotalCount: totalCount,
	}, nil
}

func (u *userService) AdminUpdateUser(targetID uint, req *requests.AdminUpdateUserRequest) (*responses.UserResponse, error) {
	user, err := u.repo.FindByID(targetID)
	if err != nil {
		return nil, configs.UserNotFound
	}

	if req.User.Username != nil {
		existingUser, err := u.repo.FindByUsername(*req.User.Username)
		if err == nil && existingUser.ID != targetID {
			return nil, configs.UsernameTaken
		}
		user.Username = *req.User.Username
	}

	if req.User.Email != nil {
		existingUser, err := u.repo.FindByEmail(*req.User.Email)
		if err == nil && existingUser.ID != targetID {
			return nil, configs.EmailTaken
		}
		user.Email = *req.User.Email
	}

	if req.User.Role != nil {
		user.Role = *req.User.Role
	}

	if err := u.repo.UpdateUser(user); err != nil {
		return nil, configs.UpdateUserFailed
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

func (u *userService) AdminDeleteUser(targetID uint) error {
	user, err := u.repo.FindByID(targetID)
	if err != nil {
		return configs.UserNotFound
	}

	if err := u.repo.DeleteUser(user); err != nil {
		return configs.DeleteUserFailed
	}
	return nil
}
