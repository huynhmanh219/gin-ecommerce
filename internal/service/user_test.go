package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"huynhmanh.com/gin/internal/model"
	"huynhmanh.com/gin/internal/util"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Create(ctx context.Context, user *model.User) (*model.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepo) GetByID(ctx context.Context, id uint) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepo) GetAll(ctx context.Context, limit, offset int) ([]model.User, int64, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]model.User), int64(args.Int(1)), args.Error(2)
}

func (m *MockUserRepo) Update(ctx context.Context, user *model.User) (*model.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepo) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockUserRepo) CreateTx(tx *gorm.DB, user *model.User) (*model.User, error) {
	args := m.Called(tx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}
func (m *MockUserRepo) UpdateTx(tx *gorm.DB, user *model.User) (*model.User, error) {
	args := m.Called(tx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}
func (m *MockUserRepo) DeleteTx(tx *gorm.DB, userID uint)  error {
	args := m.Called(tx, userID)
	if args.Get(0) == nil {
		return  args.Error(1)
	}
	return nil
}

func TestLogin_Success(t *testing.T) {
	repo := new(MockUserRepo)
	svc := NewUserService(repo)

	hash, _ := util.HashPassword("password123")
	repo.On("GetByEmail", mock.Anything,
		"huynh@example.com").Return(&model.User{
		ID:       1,
		Email:    "huynh@example.com",
		Password: hash,
	}, nil)

	token, err := svc.Login(context.Background(), "huynh@example.com", "password123")

	require.NoError(t, err)
	require.NotEmpty(t, token)

}

func TestLogin_WrongPassword(t *testing.T){
	repo := new(MockUserRepo)
	svc := NewUserService(repo)

	hash,_ := util.HashPassword("password123")
	repo.On("GetByEmail",mock.Anything,"huynh@example.com").Return(&model.User{
		ID:1,
		Email:"huynh@example.com",
		Password: hash,
	})
	
	_,err := svc.Login(context.Background(),"huynh@example.com","wrong-password")

	require.Error(t,err)
}

func TestLogin_UserNotFound(t *testing.T){
	repo := new(MockUserRepo)
	svc := NewUserService(repo)

	repo.On("GetByEmail",mock.Anything,"notfound@example.com").Return(nil,errors.New("User not found"))

	_,err := svc.Login(context.Background(),"notfound@example.com","password123")
	require.Error(t,err)
}