package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"

	"intent/kafka"
	"intent/models"
	"intent/request"
	"intent/response"
	"intent/service"
)

// UserController định nghĩa các handler cho user
type UserController struct {
	UserService   *service.UserService
	RedisClient   *redis.Client
	KafkaProducer *kafka.Producer
}

// NewUserController khởi tạo controller với UserService và RedisClient
func NewUserController(userService *service.UserService, redisClient *redis.Client, kafkaProducer *kafka.Producer) *UserController {
	return &UserController{
		UserService:   userService,
		RedisClient:   redisClient,
		KafkaProducer: kafkaProducer,
	}
}

// CreateUserHandler
func (uc *UserController) CreateUserHandler(c echo.Context) error {
	var reqs []request.UserCreateRequest
	if err := c.Bind(&reqs); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(400, "Invalid request format", err.Error()))
	}

	// Validate từng request
	for _, req := range reqs {
		if err := request.ValidateRequest(req); err != nil {
			return c.JSON(http.StatusBadRequest, response.ErrorResponse(400, "Validation failed", err.Error()))
		}
	}

	var createdUsers []models.Users
	var failedUsers []map[string]string

	requestID := uuid.New().String()

	for _, req := range reqs {
		user, err := uc.UserService.CreateUser(req)
		if err != nil {
			log.Println("Failed to create user:", err)
			failedUsers = append(failedUsers, map[string]string{
				"email": req.Email,
				"phone": req.Phone,
				"error": err.Error(),
			})
			continue
		}
		uc.pushToLog(user.ID, "CREATE")
		uc.LogUserActionToKafka("POST", "/users", req.Email, requestID)

		createdUsers = append(createdUsers, *user)
	}

	if len(createdUsers) == 0 {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(400, "All users failed to create", failedUsers))
	}

	if len(failedUsers) > 0 {
		return c.JSON(http.StatusMultiStatus, map[string]interface{}{
			"error_code": 0,
			"message":    "Some users failed to create",
			"data":       createdUsers,
			"errors":     failedUsers,
		})
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(0, "Users created successfully", createdUsers))
}

// UpdateUserHandler
func (uc *UserController) UpdateUserHandler(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(400, "Invalid ID parameter", err.Error()))
	}

	var req request.UserUpdateRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(400, "Invalid request format", err.Error()))
	}

	req.ID = id

	if err := request.ValidateRequest(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(400, "Validation failed", err.Error()))
	}

	err = uc.UserService.UpdateUser([]request.UserUpdateRequest{req})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(500, "Failed to update user", err.Error()))
	}

	uc.pushToLog(req.ID, "UPDATE")

	return c.JSON(http.StatusOK, response.SuccessResponse(0, "User updated successfully", nil))
}

// DeleteUserHandler
func (uc *UserController) DeleteUserHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(400, "Invalid user ID", nil))
	}

	err = uc.UserService.DeleteUser(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.ErrorResponse(404, "User not found", nil))
	}

	uc.pushToLog(id, "DELETE")

	return c.JSON(http.StatusOK, response.SuccessResponse(0, "User deleted successfully", nil))
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

func (uc *UserController) LogUserActionToKafka(userID, action, email, requestID string) {
	if uc.KafkaProducer == nil {
		log.Println("Kafka producer not initialized")
		return
	}

	msg := kafka.LogMessage{
		UserID:    userID,
		Action:    action,
		Email:     email,
		RequestID: requestID,
	}

	err := uc.KafkaProducer.Send(msg)
	if err != nil {
		log.Println("Failed to send log:", err)
	}
}
