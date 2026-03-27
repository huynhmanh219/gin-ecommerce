package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"huynhmanh.com/gin/internal/api/v1/handler"
	"huynhmanh.com/gin/internal/cache"
	"huynhmanh.com/gin/internal/middleware"
	"huynhmanh.com/gin/internal/repository"
	"huynhmanh.com/gin/internal/service"
)

func RegisterProductRoutes(r *gin.Engine,db *gorm.DB, redisClient *redis.Client,logger *zap.Logger){
	product_repo := repository.NewMySqlProductRepository(db)
	cacheClient := cache.NewRedisCache(redisClient)
	product_service := service.NewProductService(product_repo,cacheClient,logger)
	product_handler := handler.NewProductHandler(product_service)

	v1 := r.Group("/api/v1")
	product_group := v1.Group("/products")
	product_group.Use(middleware.AuthMiddleware())
	{
		product_group.GET("",product_handler.GetAllProduct)
		product_group.GET("/:id",product_handler.GetByID)
		product_group.POST("",product_handler.CreateProduct)
		product_group.PUT("/:id",product_handler.UpdateProduct)
		product_group.GET("/get-low-stock-product",product_handler.GetLowStockProducts)
		product_group.GET("/seach",product_handler.SearchByName)
	}
}
