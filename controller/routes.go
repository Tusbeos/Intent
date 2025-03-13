package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"intent/middleware"
	"intent/repository"
	"intent/service"
)

// RegisterUserRoutes đăng ký các route
func RegisterUserRoutes(e *echo.Echo, db *gorm.DB, redisClient *redis.Client) {
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userController := NewUserController(userService)

	usergroup := e.Group("/users")

	// Chỉ cache một số route GET
	cacheableRoutes := map[string]bool{
		"/users":     true,
		"/users/:id": true,
	}

	// Đăng ký route
	usergroup.Use(middleware.RedisCache(redisClient, cacheableRoutes))
	usergroup.POST("", userController.CreateUserHandler)
	usergroup.GET("/:id", userController.GetUserByIDHandler)
	usergroup.PUT("/:id", userController.UpdateUserHandler)
	usergroup.DELETE("/:id", userController.DeleteUserHandler)
	usergroup.GET("", userController.GetListUsersHandler)
}
