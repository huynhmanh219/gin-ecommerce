package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"huynhmanh.com/gin/internal/model"
)

type UserRepository interface {
	// regular methods
	Create(ctx context.Context, user *model.User) (*model.User,error)
	GetByID(ctx context.Context,id uint)(*model.User,error)
	GetByEmail(ctx context.Context,email string)(*model.User,error)
	GetAll(ctx context.Context,limit,offset int)([]model.User,int64,error)
	Update(ctx context.Context,user *model.User)(*model.User,error)
	Delete(ctx context.Context, id uint) error

	// transaction method 
	CreateTx(tx *gorm.DB, user *model.User)(*model.User, error)
	UpdateTx(tx *gorm.DB, user *model.User)(*model.User, error)
	DeleteTx(tx *gorm.DB,id uint)(error)

}


type MySQLUserRepository struct{
	db *gorm.DB
}

func NewMySQLUserRepository(db *gorm.DB) UserRepository{
	return &MySQLUserRepository{
		db:db,
	}
}

// transaction
func (r *MySQLUserRepository) CreateUserWithProfileTx(ctx context.Context, user *model.User,profileReq struct{Bio, Avatar string},)(*model.User,error){
	var createUser *model.User
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return fmt.Errorf("fail to create user: %w",err)
		}
		createUser = user

		profile := model.UserProFile{
			UserID: createUser.ID,
			Bio: profileReq.Bio,
			Avatar: profileReq.Avatar,
		}
		
		if err:= tx.Create(&profile).Error; err != nil {
			return fmt.Errorf("fail to create profile: %W",err)
		}

		auditLog := model.AuditLog{
			UserID: createUser.ID,
			Action: "USER_CREATED",
			Details: `{"name":"` + createUser.Name + `"}`,
			IPAddress: "127.0.0.1",
		}

		if err := tx.Create(&auditLog).Error; err !=nil{
			return fmt.Errorf("fail to create audit log : %W",err)
		}
		return nil
	})

	if err != nil{
		return nil,err
	}
	return createUser,nil
}


func (r *MySQLUserRepository) Create(ctx context.Context,user *model.User)(*model.User,error){
	if err:= r.db.WithContext(ctx).Create(user).Error; err != nil{
		return nil,err
	}
	return user,nil
}
func (r *MySQLUserRepository) GetByID(ctx context.Context,id uint)(*model.User,error){
	var user model.User

	if err:= r.db.WithContext(ctx).First(&user,id).Error; err != nil{
		if errors.Is(err,gorm.ErrRecordNotFound){
			return nil,errors.New("user not found")
		}
	}
	return &user,nil
}
func (r *MySQLUserRepository) GetByEmail(ctx context.Context,email string)(*model.User,error){
	var user model.User
	
	if err:= r.db.WithContext(ctx).Where("email = ?",email).First(&user).Error; err != nil {
		if errors.Is(err,gorm.ErrRecordNotFound){
			return nil,errors.New("user not found")
		}
		return nil,err
	}
	return &user,nil
}

func (r *MySQLUserRepository) GetAll(ctx context.Context, limit , offset int) ([]model.User,int64,error){
	var users []model.User
	var total int64

	if err:= r.db.WithContext(ctx).Model(&model.User{}).Count(&total).Error; err != nil{
		return nil,0,err
	}

	if err:= r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil,0,err
	}
	return users,total,nil
}

func (r *MySQLUserRepository) Update(ctx context.Context,user *model.User)(*model.User,error){
	if err := r.db.WithContext(ctx).Save(user).Error; err !=nil {
		return nil,err
	}
	return user, nil
}

func (r *MySQLUserRepository) Delete(ctx context.Context, id uint) error{
	if err := r.db.WithContext(ctx).Delete(&model.User{},id).Error; err != nil{
		return err
	}
	return nil
}

func (r *MySQLUserRepository) CreateTx(tx *gorm.DB, user *model.User) (*model.User,error){
	if err := tx.Create(user).Error; err != nil {
		return nil,err
	}
	return user,nil
}

func (r *MySQLUserRepository) UpdateTx(tx *gorm.DB,user *model.User)(*model.User,error){
	if err:= tx.Save(user).Error; err != nil {
		return nil,err
	}
	return user,nil
}

func (r *MySQLUserRepository) DeleteTx(tx *gorm.DB, id uint) error{
	return tx.Delete(&model.User{},id).Error
}