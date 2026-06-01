package handlers

import (
	"food_delivery/internal/configs"
	"food_delivery/internal/models/dto/requests"
	"food_delivery/internal/services"
	"food_delivery/internal/utils/constants"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body requests.RegisterRequest true "User registration data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 "Invalid request"
// @Router /user/register [post]
func (uc *UserHandler) Register(ctx *gin.Context) {
	var req requests.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Println("Error binding JSON:", err)
		ResponseError(ctx, err)
		return
	}

	user, err := uc.userService.Register(req)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	ResponseSuccess(ctx, http.StatusCreated, user)
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and receive JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body requests.LoginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 401 "Unauthorized"
// @Router /user/login [post]
func (uc *UserHandler) Login(ctx *gin.Context) {
	var req requests.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ResponseError(ctx, err)
		return
	}

	user, err := uc.userService.Login(req)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	ResponseSuccess(ctx, http.StatusOK, user)
}

// GetCurrentUser godoc
// @Summary Get current user profile
// @Description Get the profile of the authenticated user
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 "Unauthorized"
// @Router /user [get]
func (uc *UserHandler) GetCurrentUser(ctx *gin.Context) {
	userID, exists := ctx.Get(constants.UserID)
	if !exists {
		ResponseError(ctx, configs.UserNotFound)
		return
	}

	user, err := uc.userService.GetCurrentUser(userID.(uint))
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	ResponseSuccess(ctx, http.StatusOK, user)
}

// UpdateUser godoc
// @Summary Update current user profile
// @Description Update the profile of the authenticated user
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body requests.UpdateUserRequest true "Updated user data"
// @Success 200 {object} map[string]interface{}
// @Failure 401 "Unauthorized"
// @Router /user [put]
func (uc *UserHandler) UpdateUser(ctx *gin.Context) {
	userID, exists := ctx.Get(constants.UserID)
	if !exists {
		ResponseError(ctx, configs.UserNotFound)
		return
	}

	var req requests.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ResponseError(ctx, err)
		return
	}

	updatedUser, err := uc.userService.UpdateUser(userID.(uint), &req)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	ResponseSuccess(ctx, http.StatusOK, updatedUser)
}

// GetAllUsers godoc
// @Summary Get all users
// @Description Get all users with pagination (Admin only)
// @Tags Admin - Users
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {array} map[string]interface{}
// @Failure 401 "Unauthorized"
// @Failure 403 "Forbidden"
// @Router /admin/users [get]
func (uc *UserHandler) GetAllUsers(ctx *gin.Context) {
	page, limit := GetPaginationParams(ctx)

	users, err := uc.userService.GetAllUsers(page, limit)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	ResponseSuccess(ctx, http.StatusOK, users)
}

// AdminUpdateUser godoc
// @Summary Update user (Admin only)
// @Description Update any user profile (Admin only)
// @Tags Admin - Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param body body requests.AdminUpdateUserRequest true "User update data"
// @Success 200 "User updated"
// @Failure 401 "Unauthorized"
// @Failure 403 "Forbidden"
// @Router /admin/users/{id} [put]
func (uc *UserHandler) AdminUpdateUser(ctx *gin.Context) {
	targetID, err := GetIDParam(ctx)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	var req requests.AdminUpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ResponseError(ctx, err)
		return
	}

	updatedUser, err := uc.userService.AdminUpdateUser(targetID, &req)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	ResponseSuccess(ctx, http.StatusOK, updatedUser)
}

// AdminDeleteUser godoc
// @Summary Delete user (Admin only)
// @Description Delete a user (Admin only)
// @Tags Admin - Users
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 "User deleted"
// @Failure 401 "Unauthorized"
// @Failure 403 "Forbidden"
// @Router /admin/users/{id} [delete]
func (uc *UserHandler) AdminDeleteUser(ctx *gin.Context) {
	targetID, err := GetIDParam(ctx)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	if err := uc.userService.AdminDeleteUser(targetID); err != nil {
		ResponseError(ctx, err)
		return
	}

	ResponseSuccess(ctx, http.StatusOK, gin.H{"message": "User deleted successfully"})
}
