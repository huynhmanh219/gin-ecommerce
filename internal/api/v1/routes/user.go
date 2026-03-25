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

func RegisterUserRoutes(r *gin.Engine, db *gorm.DB, redisClient *redis.Client, logger *zap.Logger) {
	userRepo := repository.NewMySQLUserRepository(db)
	profileRepo := repository.NewMySQLProfileRepository(db)
	auditRepo := repository.NewMySQLAuditRepository(db)

	cacheClient := cache.NewRedisCache(redisClient)

	userService := service.NewCachedUserService(userRepo, cacheClient, logger)
	userRegisterService := service.NewUserRegistrationService(db, userRepo, profileRepo, auditRepo)
	userHandler := handler.NewUserHandler(userService)
	userRegisterHandler := handler.NewUserRegistratrionHandler(userRegisterService)

	r.GET("/api/admin/cache-metrics", middleware.CacheMetricsHandler(redisClient, logger))
	v1 := r.Group("/api/v1")

	v1.POST("/users", userHandler.CreateUser)
	v1.POST("/auth/login", userHandler.Login)

	userGroup := v1.Group("/users")
	userGroup.Use(middleware.AuthMiddleware())
	{
		userGroup.GET("", userHandler.GetUsers)
		userGroup.GET("/:id", userHandler.GetUser)
		userGroup.PUT("/:id", userHandler.UpdateUser)
		userGroup.DELETE("/:id", userHandler.DeleteUser)
		userGroup.POST("/register-full", userRegisterHandler.RegisterWithProfile)
	}

}
