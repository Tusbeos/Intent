package controller

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"Intent/repository"
	"Intent/service"
)

// RegisterUserRoutes đăng ký các route liên quan đến user
func RegisterUserRoutes(e *echo.Echo, db *gorm.DB) {
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userController := NewUserController(userService)

	e.POST("/users", userController.CreateUserHandler)
	e.GET("/users/:id", userController.GetUserByIDHandler)
	e.PUT("/users/:id", userController.UpdateUserHandler)
	e.DELETE("/users/:id", userController.DeleteUserHandler)
	e.GET("/users", userController.GetListUsersHandler)
}
