package handlers

import (
	"food_delivery/internal/configs"
	"food_delivery/internal/models/dto/requests"
	"food_delivery/internal/services"
	"food_delivery/internal/utils/constants"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

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

func (uc *UserHandler) GetAllUsers(ctx *gin.Context) {
	users, err := uc.userService.GetAllUsers()
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	ResponseSuccess(ctx, http.StatusOK, users)
}

func (uc *UserHandler) AdminUpdateUser(ctx *gin.Context) {
	targetID, err := getIDParam(ctx)
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

func (uc *UserHandler) AdminDeleteUser(ctx *gin.Context) {
	targetID, err := getIDParam(ctx)
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

func getIDParam(ctx *gin.Context) (uint, error) {
	idStr := ctx.Param(constants.IDParam)
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
