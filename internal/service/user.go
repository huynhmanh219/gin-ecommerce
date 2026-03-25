package service

import (
	"context"
	"errors"


	"huynhmanh.com/gin/internal/model"
	"huynhmanh.com/gin/internal/repository"
	"huynhmanh.com/gin/internal/util"
)

type UserService interface {
	Register(ctx context.Context, name, email, password string)(*model.User, error)

	GetUserByID(ctx context.Context,id uint)(*model.User,error)

	GetAllUsers(ctx context.Context,limit,offset int)([]model.User,int64,error)

	UpdateUser(ctx context.Context,id uint,name,email string)(*model.User,error)

	DeleteUser(ctx context.Context,id uint) error

	Login(ctx context.Context,email,password string)(string,error)
}

type userService struct{
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService{
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) Register(ctx context.Context,name,email,password string)(*model.User,error){
	if len(name) ==0 || len(email) == 0 || len(password) < 6{
		return nil, errors.New("invalid input: name/email required, password must greater than 6 characters")
	}
	existingUser,_ := s.userRepo.GetByEmail(ctx,email)
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	hashedPassword,err := util.HashPassword(password)
	if err != nil {
		return nil,err
	}

	user:= &model.User{
		Name: name,
		Email: email,
		Password: hashedPassword,
		IsActive: true,
	}
	createUser,err := s.userRepo.Create(ctx,user)
	if err != nil{
		return nil,err
	}
	return createUser,nil
}
func (s *userService) GetUserByID(ctx context.Context,id uint)(*model.User,error){
	user, err := s.userRepo.GetByID(ctx,id)
	if err != nil {
		return nil, err
	}

	user.Password = ""
	return user,nil
}

func (s *userService) GetAllUsers(ctx context.Context, limit, offset int)([]model.User,int64,error){
	users,total, err := s.userRepo.GetAll(ctx,limit,offset)
	if err != nil {
		return nil,0,err
	}

	for i := range users{
		users[i].Password= ""
	}
	return users,total,nil
}
func (s *userService) UpdateUser(ctx context.Context,id uint,name,email string) (*model.User,error){
	user,err := s.userRepo.GetByID(ctx,id)
	if err != nil {
		return nil, err
	}

	if len(name) > 0 {
		user.Name = name
	}

	if len(email) > 0 {
		existngUser,_ := s.userRepo.GetByEmail(ctx,email)
		if existngUser != nil && existngUser.ID != id{
			return nil,errors.New("email already taken")
		}
		user.Email = email
	}

	updateUser,err :=s.userRepo.Update(ctx,user)
	if err != nil {
		return nil,err
	}

	updateUser.Password = ""
	return updateUser,nil
}

func (s *userService) DeleteUser(ctx context.Context,id uint)error{
	_,err := s.userRepo.GetByID(ctx,id)
	if err != nil{
		return err
	}
	return s.userRepo.Delete(ctx,id)
}

func (s *userService) Login(ctx context.Context,email,password string)(string, error){
	user,err := s.userRepo.GetByEmail(ctx,email)
	if err != nil{
		return "",errors.New("Invlid credentials")
	}
	
	if !util.VerifyPassword(user.Password,password){
		return "",errors.New("password not correct")
	}
	token,err := util.GenerateToken(user.ID,user.Email,24)
	if err != nil{
		return "",err
	}
	return token,nil
}