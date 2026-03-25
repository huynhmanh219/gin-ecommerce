package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"huynhmanh.com/gin/internal/model"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Register(ctx context.Context, name, email, password string) (*model.User, error) {
	args := m.Called(ctx, name, email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) Login(ctx context.Context, email, password string) (string, error) {
	args := m.Called(ctx, email, password)
	return args.String(0), args.Error(1)
}

func (m *MockUserService) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) GetAllUsers(ctx context.Context, limit, offset int) ([]model.User, int64, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]model.User), int64(args.Int(1)), args.Error(2)
}

func (m *MockUserService) UpdateUser(ctx context.Context, id uint, name, email string) (*model.User, error) {
	args := m.Called(ctx, id, name, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) DeleteUser(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateUser_BadRequest_InvalidJSON(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	r := gin.New()
	r.POST("/api/v1/users", handler.CreateUser)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/users",
		strings.NewReader(`{"email":"test@example.com"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.NotNil(t, response["error"])
}


func TestCreateUser_ValidationError_InvalidEmail(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	r := gin.New()
	r.POST("/api/v1/users", handler.CreateUser)

	// Email không đúng format
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users",
		strings.NewReader(`{"name":"Test User","email":"invalid-email","password":"password123"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.Equal(t, "Validation_error", response["code"])
	require.NotNil(t, response["errors"])
}


func TestCreateUser_ValidationError_ShortPassword(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	r := gin.New()
	r.POST("/api/v1/users", handler.CreateUser)

	
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users",
		strings.NewReader(`{"name":"Test User","email":"test@example.com","password":"123"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.Equal(t, "Validation_error", response["code"])
}


func TestCreateUser_ValidationError_ShortName(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	r := gin.New()
	r.POST("/api/v1/users", handler.CreateUser)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/users",
		strings.NewReader(`{"name":"A","email":"test@example.com","password":"password123"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.Equal(t, "Validation_error", response["code"])
}

// Test case: Service error - email đã tồn tại
func TestCreateUser_Conflict_DuplicateEmail(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	// Mock service trả về error (email duplicate)
	mockService.On("Register", mock.Anything, "Test User", "duplicate@example.com", "password123").
		Return(nil, errors.New("email already exists"))

	r := gin.New()
	r.POST("/api/v1/users", handler.CreateUser)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/users",
		strings.NewReader(`{"name":"Test User","email":"duplicate@example.com","password":"password123"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.Equal(t, "email already exists", response["error"])
}

// Test case: Success - tạo user thành công
func TestCreateUser_Success(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	createdAt := time.Now()
	newUser := &model.User{
		ID:        1,
		Name:      "Test User",
		Email:     "test@example.com",
		IsActive:  true,
		CreatedAt: createdAt,
	}

	mockService.On("Register", mock.Anything, "Test User", "test@example.com", "password123").
		Return(newUser, nil)

	r := gin.New()
	r.POST("/api/v1/users", handler.CreateUser)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/users",
		strings.NewReader(`{"name":"Test User","email":"test@example.com","password":"password123"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].(map[string]interface{})
	require.Equal(t, float64(1), data["id"])
	require.Equal(t, "Test User", data["name"])
	require.Equal(t, "test@example.com", data["email"])
	require.Equal(t, true, data["is_active"])
}
