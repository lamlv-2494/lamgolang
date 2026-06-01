package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"food_delivery/internal/configs"
	"food_delivery/internal/models/dto/requests"
	"food_delivery/internal/models/dto/responses"
	"food_delivery/internal/utils/constants"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// helper: create a test gin context with a JSON body
func newTestCtx(body any) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	if body != nil {
		b, _ := json.Marshal(body)
		ctx.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(b))
		ctx.Request.Header.Set("Content-Type", "application/json")
	} else {
		ctx.Request = httptest.NewRequest(http.MethodPost, "/", nil)
	}
	return ctx, w
}

// ─── Register ─────────────────────────────────────────────────────────────────

func TestUserHandler_Register_Success(t *testing.T) {
	svc := &mockUserService{
		registerFn: func(req requests.RegisterRequest) (*responses.UserResponse, error) {
			return &responses.UserResponse{User: responses.UserData{Email: "a@b.com"}}, nil
		},
	}
	h := NewUserHandler(svc)
	ctx, w := newTestCtx(map[string]interface{}{"user": map[string]string{"email": "a@b.com", "username": "alice", "password": "secret"}})
	h.Register(ctx)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestUserHandler_Register_BindError(t *testing.T) {
	svc := &mockUserService{}
	h := NewUserHandler(svc)
	ctx, w := newTestCtx(nil) // no body → bind error
	ctx.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("invalid json{"))
	ctx.Request.Header.Set("Content-Type", "application/json")
	h.Register(ctx)
	assert.NotEqual(t, http.StatusCreated, w.Code)
}

func TestUserHandler_Register_ServiceError(t *testing.T) {
	svc := &mockUserService{
		registerFn: func(req requests.RegisterRequest) (*responses.UserResponse, error) {
			return nil, configs.EmailOrUserTaken
		},
	}
	h := NewUserHandler(svc)
	ctx, w := newTestCtx(map[string]interface{}{"user": map[string]string{"email": "a@b.com", "username": "alice", "password": "secret"}})
	h.Register(ctx)
	assert.NotEqual(t, http.StatusCreated, w.Code)
}

// ─── Login ────────────────────────────────────────────────────────────────────

func TestUserHandler_Login_Success(t *testing.T) {
	svc := &mockUserService{
		loginFn: func(req requests.LoginRequest) (*responses.UserResponse, error) {
			return &responses.UserResponse{User: responses.UserData{Email: "a@b.com"}}, nil
		},
	}
	h := NewUserHandler(svc)
	ctx, w := newTestCtx(map[string]interface{}{"email": "a@b.com", "password": "secret"})
	h.Login(ctx)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserHandler_Login_BindError(t *testing.T) {
	svc := &mockUserService{}
	h := NewUserHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("{bad json"))
	ctx.Request.Header.Set("Content-Type", "application/json")
	h.Login(ctx)
	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestUserHandler_Login_ServiceError(t *testing.T) {
	svc := &mockUserService{
		loginFn: func(req requests.LoginRequest) (*responses.UserResponse, error) {
			return nil, configs.UserNotFound
		},
	}
	h := NewUserHandler(svc)
	ctx, w := newTestCtx(map[string]interface{}{"email": "x@y.com", "password": "wrong"})
	h.Login(ctx)
	assert.NotEqual(t, http.StatusOK, w.Code)
}

// ─── GetCurrentUser ───────────────────────────────────────────────────────────

func TestUserHandler_GetCurrentUser_Success(t *testing.T) {
	svc := &mockUserService{
		getCurrentUserFn: func(id uint) (*responses.UserResponse, error) {
			return &responses.UserResponse{User: responses.UserData{ID: id}}, nil
		},
	}
	h := NewUserHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	ctx.Set(constants.UserID, uint(1))
	h.GetCurrentUser(ctx)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserHandler_GetCurrentUser_NoUserID(t *testing.T) {
	svc := &mockUserService{}
	h := NewUserHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	h.GetCurrentUser(ctx)
	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestUserHandler_GetCurrentUser_ServiceError(t *testing.T) {
	svc := &mockUserService{
		getCurrentUserFn: func(id uint) (*responses.UserResponse, error) {
			return nil, configs.UserNotFound
		},
	}
	h := NewUserHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	ctx.Set(constants.UserID, uint(99))
	h.GetCurrentUser(ctx)
	assert.NotEqual(t, http.StatusOK, w.Code)
}

// ─── UpdateUser ───────────────────────────────────────────────────────────────

func TestUserHandler_UpdateUser_Success(t *testing.T) {
	svc := &mockUserService{
		updateUserFn: func(id uint, req *requests.UpdateUserRequest) (*responses.UserResponse, error) {
			return &responses.UserResponse{}, nil
		},
	}
	h := NewUserHandler(svc)
	ctx, w := newTestCtx(map[string]interface{}{"user": map[string]string{"username": "newname"}})
	ctx.Set(constants.UserID, uint(1))
	h.UpdateUser(ctx)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserHandler_UpdateUser_NoUserID(t *testing.T) {
	svc := &mockUserService{}
	h := NewUserHandler(svc)
	ctx, w := newTestCtx(map[string]interface{}{"user": map[string]string{"username": "newname"}})
	h.UpdateUser(ctx)
	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestUserHandler_UpdateUser_BindError(t *testing.T) {
	svc := &mockUserService{}
	h := NewUserHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodPut, "/", bytes.NewBufferString("{bad"))
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Set(constants.UserID, uint(1))
	h.UpdateUser(ctx)
	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestUserHandler_UpdateUser_ServiceError(t *testing.T) {
	svc := &mockUserService{
		updateUserFn: func(id uint, req *requests.UpdateUserRequest) (*responses.UserResponse, error) {
			return nil, configs.UserNotFound
		},
	}
	h := NewUserHandler(svc)
	ctx, w := newTestCtx(map[string]interface{}{"user": map[string]string{"username": "x"}})
	ctx.Set(constants.UserID, uint(1))
	h.UpdateUser(ctx)
	assert.NotEqual(t, http.StatusOK, w.Code)
}

// ─── GetAllUsers ──────────────────────────────────────────────────────────────

func TestUserHandler_GetAllUsers_Success(t *testing.T) {
	svc := &mockUserService{
		getAllUsersFn: func(p, l int) (*responses.ListResponse[*responses.UserData], error) {
			return &responses.ListResponse[*responses.UserData]{}, nil
		},
	}
	h := NewUserHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/?page=1&limit=10", nil)
	h.GetAllUsers(ctx)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserHandler_GetAllUsers_ServiceError(t *testing.T) {
	svc := &mockUserService{
		getAllUsersFn: func(p, l int) (*responses.ListResponse[*responses.UserData], error) {
			return nil, errors.New("db error")
		},
	}
	h := NewUserHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	h.GetAllUsers(ctx)
	assert.NotEqual(t, http.StatusOK, w.Code)
}

// ─── AdminUpdateUser ──────────────────────────────────────────────────────────

func TestUserHandler_AdminUpdateUser_Success(t *testing.T) {
	svc := &mockUserService{
		adminUpdateUserFn: func(id uint, req *requests.AdminUpdateUserRequest) (*responses.UserResponse, error) {
			return &responses.UserResponse{}, nil
		},
	}
	h := NewUserHandler(svc)
	ctx, w := newTestCtx(map[string]interface{}{"user": map[string]string{"username": "admin_user"}})
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	h.AdminUpdateUser(ctx)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserHandler_AdminUpdateUser_InvalidID(t *testing.T) {
	svc := &mockUserService{}
	h := NewUserHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodPut, "/", nil)
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "abc"}}
	h.AdminUpdateUser(ctx)
	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestUserHandler_AdminUpdateUser_BindError(t *testing.T) {
	svc := &mockUserService{}
	h := NewUserHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodPut, "/", bytes.NewBufferString("{bad"))
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	h.AdminUpdateUser(ctx)
	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestUserHandler_AdminUpdateUser_ServiceError(t *testing.T) {
	svc := &mockUserService{
		adminUpdateUserFn: func(id uint, req *requests.AdminUpdateUserRequest) (*responses.UserResponse, error) {
			return nil, configs.UserNotFound
		},
	}
	h := NewUserHandler(svc)
	ctx, w := newTestCtx(map[string]interface{}{"user": map[string]string{"username": "x"}})
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	h.AdminUpdateUser(ctx)
	assert.NotEqual(t, http.StatusOK, w.Code)
}

// ─── AdminDeleteUser ──────────────────────────────────────────────────────────

func TestUserHandler_AdminDeleteUser_Success(t *testing.T) {
	svc := &mockUserService{
		adminDeleteUserFn: func(id uint) error { return nil },
	}
	h := NewUserHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodDelete, "/", nil)
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	h.AdminDeleteUser(ctx)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserHandler_AdminDeleteUser_InvalidID(t *testing.T) {
	svc := &mockUserService{}
	h := NewUserHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodDelete, "/", nil)
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "abc"}}
	h.AdminDeleteUser(ctx)
	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestUserHandler_AdminDeleteUser_ServiceError(t *testing.T) {
	svc := &mockUserService{
		adminDeleteUserFn: func(id uint) error { return configs.UserNotFound },
	}
	h := NewUserHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodDelete, "/", nil)
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	h.AdminDeleteUser(ctx)
	assert.NotEqual(t, http.StatusOK, w.Code)
}
