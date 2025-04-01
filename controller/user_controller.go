package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"

	"intent/models"
	"intent/request"
	"intent/response"
	"intent/service"
)

// UserController định nghĩa các handler cho user
type UserController struct {
	UserService *service.UserService
	RedisClient *redis.Client
}

// NewUserController khởi tạo controller với UserService và RedisClient
func NewUserController(userService *service.UserService, redisClient *redis.Client) *UserController {
	return &UserController{
		UserService: userService,
		RedisClient: redisClient,
	}
}

// CreateUserHandler
func (uc *UserController) CreateUserHandler(c echo.Context) error {
	var reqs []request.UserCreateRequest
	if err := c.Bind(&reqs); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(400, "Invalid request format", err.Error()))
	}

	// Validate từng request trong danh sách
	for _, req := range reqs {
		if err := request.ValidateRequest(req); err != nil {
			return c.JSON(http.StatusBadRequest, response.ErrorResponse(400, "Validation failed", err.Error()))
		}
	}

	var createdUsers []models.Users
	for _, req := range reqs {
		user, err := uc.UserService.CreateUser(req)
		if err != nil {
			log.Println("Failed to create user:", err)
			continue
		}

		uc.pushToLog(user.ID, "CREATE")
		createdUsers = append(createdUsers, *user)
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(0, "Users created successfully", createdUsers))
}

// UpdateUserHandler
func (uc *UserController) UpdateUserHandler(c echo.Context) error {
	var reqs []request.UserUpdateRequest
	if err := c.Bind(&reqs); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(400, "Invalid request format", err.Error()))
	}

	// Validate request
	for _, req := range reqs {
		if err := request.ValidateRequest(req); err != nil {
			return c.JSON(http.StatusBadRequest, response.ErrorResponse(400, "Validation failed", err.Error()))
		}
	}

	err := uc.UserService.UpdateUser(reqs)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(500, "Failed to update users", err.Error()))
	}

	// Gửi event cập nhật
	for _, req := range reqs {
		uc.pushToLog(req.ID, "UPDATE")
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(0, "Users updated successfully", nil))
}

// GetUserByIDHandler
func (uc *UserController) GetUserByIDHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(400, "Invalid user ID", nil))
	}

	user, err := uc.UserService.GetUserByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.ErrorResponse(404, "User not found", nil))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(0, "User retrieved successfully", user))
}

// DeleteUserHandler - Giữ nguyên
func (uc *UserController) DeleteUserHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(400, "Invalid user ID", nil))
	}

	err = uc.UserService.DeleteUser(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.ErrorResponse(404, "User not found", nil))
	}

	// Bắn message vào Redis queue
	uc.pushToLog(id, "DELETE")

	return c.JSON(http.StatusOK, response.SuccessResponse(0, "User deleted successfully", nil))
}

// GetListUsersHandler
func (uc *UserController) GetListUsersHandler(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	status := c.QueryParam("status")
	gender := c.QueryParam("gender")

	req := request.GetListUsersRequest{
		Page:   page,
		Limit:  limit,
		Status: status,
		Gender: gender,
	}

	users, meta, err := uc.UserService.GetListUsers(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(500, "Failed to get users", err.Error()))
	}
	return c.JSON(http.StatusOK, response.SuccessResponseWithMeta(0, "Users retrieved successfully", users, meta))
}
