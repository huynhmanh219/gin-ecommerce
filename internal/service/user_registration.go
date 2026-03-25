package service

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"huynhmanh.com/gin/internal/dto"
	"huynhmanh.com/gin/internal/model"
	"huynhmanh.com/gin/internal/repository"
)

type UserRegistrationService interface {
	RegisterWithProfile(
		ctx context.Context,
		req dto.CreateUserRequest,
		bio,avatar string,
	)(*model.User,error)
}
type userRegistrationService struct{
	db *gorm.DB
	userRepo repository.UserRepository
	profileRepo repository.ProfileRepository
	auditRepo repository.AuditRepository
}

func NewUserRegistrationService(
	db *gorm.DB,
	userRepo repository.UserRepository,
	profileRepo repository.ProfileRepository,
	auditRepo repository.AuditRepository,
)UserRegistrationService {
	return &userRegistrationService{
		db:        db,
        userRepo:  userRepo,
        profileRepo: profileRepo,
        auditRepo: auditRepo,
	}
}

func (s *userRegistrationService) RegisterWithProfile(
	ctx context.Context,
	req dto.CreateUserRequest,
	bio, avatar string,
)(*model.User,error){
	var createdUser *model.User

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		hashedPassword,err := bcrypt.GenerateFromPassword([]byte(req.Password),10)
		if err != nil {return nil}

		user := &model.User{
			Name: req.Name,
			Email: req.Email,
			Password: string(hashedPassword),
			IsActive: true,
		}

		createdUser, err := s.userRepo.CreateTx(tx,user)
		if err != nil {
			return fmt.Errorf("failto create user: %W",err)
		}

		profile := &model.UserProFile{
			UserID: createdUser.ID,
			Bio: bio,
			Avatar: avatar,
		}
		if _,err := s.profileRepo.CreateTx(tx,profile); err != nil{
			return fmt.Errorf("Failed to create profile: %W", err)
		}
		
		auditLog := &model.AuditLog{
            UserID:    createdUser.ID,
            Action:    "USER_REGISTERED",
            Details:   fmt.Sprintf(`{"name":"%s","email":"%s"}`, createdUser.Name, createdUser.Email),
        }
        
        if _,err := s.auditRepo.CreateTx(tx, auditLog); err != nil {
            return fmt.Errorf("failed to write audit log: %w", err)
        }
		return nil
	})
	if err != nil {
		return nil,err
	}
	return createdUser,nil
}