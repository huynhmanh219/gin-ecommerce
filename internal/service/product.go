package service

import (
	"context"
	"encoding/json"
	"time"

	"go.uber.org/zap"
	"huynhmanh.com/gin/internal/cache"
	"huynhmanh.com/gin/internal/dto"
	"huynhmanh.com/gin/internal/model"
	"huynhmanh.com/gin/internal/repository"
)

type ProductService interface {
	GetAllProduct(ctx context.Context,limit,offset uint) ([]model.Product,int64,error)
	GetByID(ctx context.Context,id uint) (*model.Product,error)
	CreateProduct(ctx context.Context, input dto.CreateProductRequest) (*model.Product,error)
	UpdateProduct(ctx context.Context,id uint , input dto.UpdateProductRequest)(*model.Product,error)
	DeleteProduct(ctx context.Context,id uint) error
	SearchByName(ctx context.Context,name string) (*model.Product,error)
	GetLowStockProducts(ctx context.Context)([]model.Product,error)
}


const (
	ProductCacheTTL     = 5 * time.Minute
	ProductListCacheTTL = 1 * time.Minute
)
type productService struct {
	productRepo repository.ProductRepository
	cacheConn  cache.CacheClient
	logger *zap.Logger
}

func NewProductService(
	repo repository.ProductRepository,
	cacheCoon cache.CacheClient,
	logger *zap.Logger,
) ProductService{
	return &productService{
		productRepo: repo,
		cacheConn: cacheCoon,
		logger: logger,
	}
}

func (s *productService) GetAllProduct(ctx context.Context,limit,offset uint)([]model.Product,int64,error){
	data,total,err := s.productRepo.GetAllProduct(ctx,limit,offset)
	if err != nil {
		return nil,0,err
	}

	return data,total,nil
}

func (s *productService) GetByID(ctx context.Context,id uint)(*model.Product,error){
	key:= GetProductCacheKey(id)
	
	cache,err := s.cacheConn.Get(ctx,key)
	if err == nil && cache != ""{
		s.logger.Info("cache_hit",zap.String("key",key))
		s.cacheConn.Increment(ctx,"cache:hits")

		var product model.Product
		if err := json.Unmarshal([]byte(cache),&product); err == nil {
			return &product,nil
		}
	}
	s.logger.Info("cache_miss",zap.String("key",key))
	s.cacheConn.Increment(ctx,"cache:misses")

	data,err := s.productRepo.GetByID(ctx,id)
	if err != nil{
		return nil,err
	}
	if data,err := json.Marshal(data); err == nil {
		_ = s.cacheConn.Set(ctx,key,string(data),UserCacheTTL)
	}
	return data,nil
}

func (s *productService) CreateProduct(ctx context.Context,input dto.CreateProductRequest)(*model.Product,error){
	var product *model.Product
	product = &model.Product{
		Name: input.Name,
		Description: input.Description,
		Price: input.Price,	
		Stock: input.Stock,
		Category: input.Category,
		ImageURL: input.ImageURL,
		IsActive: input.IsActive,	
	}
	if _,err := s.productRepo.Create(ctx,product); err !=nil{
		return nil,err
	}
	return product,nil
}

func (s *productService) UpdateProduct(ctx context.Context,id uint, input dto.UpdateProductRequest)(*model.Product,error){
	data,err := s.productRepo.GetByID(ctx,id)
	if err != nil {
		return nil,err
	}
	data.Name = input.Name
	data.Description = input.Description
	data.Price = input.Price
	data.Stock = input.Stock
	data.Category = input.Category
	data.ImageURL = input.ImageURL
	data.IsActive = input.IsActive
	
	if _,err := s.productRepo.Update(ctx,data);err !=nil {
		return nil,err
	}

	key:= GetProductCacheKey(id)
	if err := s.cacheConn.Delete(ctx,key); err != nil{
		s.logger.Warn("cache_invalidtion_failed",zap.String("key",key))
	}
	return data,nil
}

func (s *productService) DeleteProduct(ctx context.Context,id uint) error{
	data,err := s.productRepo.GetByID(ctx,id)
	if err != nil {
		return err
	}
	data.IsActive = false
	if _,err := s.productRepo.Update(ctx,data); err != nil{
		return err
	}
	key := GetProductCacheKey(id)
	_ = s.cacheConn.Delete(ctx,key)

	return nil
}

func (s *productService) SearchByName(ctx context.Context, name string) (*model.Product,error){
	data,err := s.productRepo.SearchByName(ctx,name)
	if err != nil {
		return nil,err
	}
	return data,nil
}

func (s *productService) GetLowStockProducts(ctx context.Context)([]model.Product,error){
	data,err := s.productRepo.GetLowStockProducts(ctx)
	if err != nil {return nil,err}
	return data,nil
}