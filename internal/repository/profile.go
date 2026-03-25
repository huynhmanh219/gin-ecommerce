package repository

import (
	"gorm.io/gorm"
	"huynhmanh.com/gin/internal/model"
)

type ProfileRepository interface {
	Create(profile *model.UserProFile)(*model.UserProFile,error)
	Update(profile *model.UserProFile)(*model.UserProFile,error)
	GetByUserID(userID uint)(*model.UserProFile,error)

	CreateTx(tx *gorm.DB, profile *model.UserProFile)(*model.UserProFile,error)
	UpdateTx(tx *gorm.DB,profile *model.UserProFile)(*model.UserProFile,error)
}

type MySQLProfileRepository struct{
	db *gorm.DB
}

func NewMySQLProfileRepository(db *gorm.DB) ProfileRepository{
	return &MySQLProfileRepository{db:db}
}

func (r *MySQLProfileRepository) CreateTx(tx *gorm.DB,profile *model.UserProFile)(*model.UserProFile,error){
	if err := tx.Create(profile).Error; err != nil {
		return nil,err
	}
	return profile,nil
}

func(r *MySQLProfileRepository) UpdateTx(tx *gorm.DB,profile *model.UserProFile)(*model.UserProFile,error){
	if err:= tx.Save(profile).Error; err != nil {
		return nil,err
	}
	return profile,nil
}

func (r *MySQLProfileRepository) Create(profile *model.UserProFile)(*model.UserProFile,error){
	if err := r.db.Create(profile).Error; err != nil {
		return nil,err
	}
	return profile,nil
}
func (r *MySQLProfileRepository) Update(profile *model.UserProFile)(*model.UserProFile,error){
	if err:= r.db.Save(profile).Error; err !=nil {
		return nil,err
	}
	return profile,nil
}

func (r *MySQLProfileRepository) GetByUserID(userID uint)(*model.UserProFile,error){
	var data model.UserProFile
	if err:= r.db.Where("userID = ?",userID).First(&data).Error; err != nil{
		return nil,err
	}
	return &data,nil
}

