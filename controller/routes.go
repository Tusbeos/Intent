package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"intent/repository"
	"intent/service"
)

// RegisterUserRoutes đăng ký các route
func RegisterUserRoutes(e *echo.Echo, db *gorm.DB, redisClient *redis.Client) {
	userRepo := repository.NewUserRepository(db, redisClient)
	userService := service.NewUserService(userRepo)
	userController := NewUserController(userService)

	usergroup := e.Group("/users")
	usergroup.POST("", userController.CreateUserHandler)
	usergroup.GET("/:id", userController.GetUserByIDHandler)
	usergroup.PUT("/:id", userController.UpdateUserHandler)
	usergroup.DELETE("/:id", userController.DeleteUserHandler)
	usergroup.GET("", userController.GetListUsersHandler)
}
