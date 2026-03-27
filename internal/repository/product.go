package repository

import (
	"context"

	"gorm.io/gorm"
	"huynhmanh.com/gin/internal/model"
)

type ProductRepository interface {
	Create(ctx context.Context, product *model.Product)(*model.Product,error)
	GetAllProduct(ctx context.Context,limit,offset uint) ([]model.Product,int64,error)
	GetByID(ctx context.Context,id uint) (*model.Product,error)
	Update(ctx context.Context,product *model.Product)(*model.Product,error)
	Delete(ctx context.Context, id uint) error
	SearchByName(ctx context.Context,name string)(*model.Product,error)
	GetLowStockProducts(ctx context.Context) ([]model.Product,error)
}
type MySqlProductRepository struct{
	db *gorm.DB
}

func NewMySqlProductRepository(db *gorm.DB) ProductRepository{
	return &MySqlProductRepository{
		db:db,
	}
}

func (r *MySqlProductRepository) Create(ctx context.Context,product *model.Product)(*model.Product,error){
	if err:= r.db.WithContext(ctx).Create(product).Error; err != nil{
		return nil,err
	}
	return product,nil
}

func (r *MySqlProductRepository) Update(ctx context.Context,product *model.Product)(*model.Product,error){
	if err:= r.db.WithContext(ctx).Save(product).Error; err != nil{
		return nil,err
	}
	return product,nil
}

func (r *MySqlProductRepository) Delete(ctx context.Context,id uint)(error){
	if err:= r.db.WithContext(ctx).Delete(id).Error; err != nil{
		return err
	}
	return nil
}
func (r *MySqlProductRepository) GetAllProduct(ctx context.Context,limit,offset uint)([]model.Product,int64,error){
	var products []model.Product
	var total int64
	if err:= r.db.WithContext(ctx).Count(&total).Error; err != nil {
		return nil,0,err
	}
	if err:= r.db.WithContext(ctx).Limit(int(limit)).Offset(int(offset)).Find(&products).Error; err != nil {
		return nil,0,err
	}
	return products,total,nil
}

func (r *MySqlProductRepository) GetByID(ctx context.Context,id uint)(*model.Product,error){
	var products model.Product
	
	if err:= r.db.WithContext(ctx).First(&products,id).Error; err != nil {
		return nil,err
	}
	return &products,nil
}
func (r *MySqlProductRepository) SearchByName(ctx context.Context,name string) (*model.Product,error){
	var data *model.Product
	if err := r.db.WithContext(ctx).Where("name = ?",name).First(&data).Error; err != nil{
		return nil,err
	}
	return data,nil
}

func (r *MySqlProductRepository) GetLowStockProducts(ctx context.Context) ([]model.Product,error) {
	var products []model.Product
	err := r.db.WithContext(ctx).Where("stock < ?",10).Find(&products).Error
	if err != nil {
		return nil,err
	}
	return products,nil
}