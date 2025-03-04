package controller

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"Intent/request"
	"Intent/response"
	"Intent/service"
)

// UserController định nghĩa các handler cho user
type UserController struct {
	UserService *service.UserService
}

// NewUserController khởi tạo controller với UserService
func NewUserController(userService *service.UserService) *UserController {
	return &UserController{UserService: userService}
}

// CreateUserHandler xử lý tạo user mới
func (uc *UserController) CreateUserHandler(c echo.Context) error {
	var req request.UserCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(400, "Invalid request", err.Error()))
	}

	user, err := uc.UserService.CreateUser(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(500, "Failed to create user", err.Error()))
	}

	return c.JSON(http.StatusCreated, response.SuccessResponse(0, "User created successfully", user))
}

// GetUserByIDHandler xử lý lấy user theo ID
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

// UpdateUserHandler xử lý cập nhật user
func (uc *UserController) UpdateUserHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(400, "Invalid user ID", nil))
	}

	var req request.UserUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(400, "Invalid request", err.Error()))
	}

	err = uc.UserService.UpdateUser(id, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(500, "Failed to update user", err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(0, "User updated successfully", nil))
}

// DeleteUserHandler xử lý xóa user
func (uc *UserController) DeleteUserHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(400, "Invalid user ID", nil))
	}

	err = uc.UserService.DeleteUser(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.ErrorResponse(404, "User not found", nil))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(0, "User deleted successfully", nil))
}

// GetListUsersHandler xử lý lấy danh sách user
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
