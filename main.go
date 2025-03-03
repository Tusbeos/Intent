package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	config "Intent/config"
	models "Intent/models"
	request "Intent/request"
)

var cfg *config.Config

func connectDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		viper.GetString("database.user"),
		viper.GetString("database.password"),
		viper.GetString("database.host"),
		viper.GetInt("database.port"),
		viper.GetString("database.dbname"),
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	log.Println("Successfully connected to MySQL!")
	return db, nil
}

func main() {
	cfg = config.LoadConfig()
	db, err := connectDB()
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	// Khởi tạo router
	http.HandleFunc("POST /users", func(w http.ResponseWriter, r *http.Request) {
		AddUser(w, r, db)
	})
	http.HandleFunc("PUT /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		UpdateUser(w, r, db)
	})
	http.HandleFunc("DELETE /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		DeleteUser(w, id, db)
	})
	http.HandleFunc("GET /users/", func(w http.ResponseWriter, r *http.Request) {
		GetListUsers(w, r, db)
	})
	http.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		GetUserByID(w, id, db)
	})

	log.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", nil)
}

// API: Thêm user
func AddUser(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var req request.UserCreateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		jsonResponse(w, 400, "Invalid JSON data", nil)
		return
	}

	// Validate request
	if errMsg := request.ValidateRequest(&req); errMsg != nil {
		log.Printf("Validation error: %v", errMsg)
		jsonResponse(w, 400, errMsg.Error(), nil)
		return
	}

	user := models.Users{
		Name:     req.Name,
		Password: req.Password,
		Email:    req.Email,
		Phone:    req.Phone,
		Gender:   req.Gender,
		Status:   req.Status,
	}
	if err := db.Create(&user).Error; err != nil {
		log.Printf("Error adding user: %v", err)
		jsonResponse(w, 500, "Error adding user", nil)
		return
	}

	jsonResponse(w, 0, "User added successfully!", nil)
}

// API: Cập nhật user
func UpdateUser(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	id := r.PathValue("id")
	if id == "" {
		log.Println("Missing user ID")
		jsonResponse(w, 400, "Missing user ID", nil)
		return
	}

	// Decode dữ liệu từ request body
	var req request.UserUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		jsonResponse(w, 400, "Invalid JSON data", nil)
		return
	}

	// Validate request
	if errMsg := request.ValidateRequest(&req); errMsg != nil {
		log.Printf("Validation error: %v", errMsg)
		jsonResponse(w, 400, errMsg.Error(), nil)
		return
	}

	// Kiểm tra user có tồn tại không
	var existingUser models.Users
	if err := db.First(&existingUser, id).Error; err != nil {
		log.Printf("User not found: %v", err)
		jsonResponse(w, 404, "User not found", nil)
		return
	}

	// Cập nhật dữ liệu user
	if err := db.Model(&existingUser).Updates(models.Users{
		Name:     req.Name,
		Password: req.Password,
		Email:    req.Email,
		Phone:    req.Phone,
		Gender:   req.Gender,
		Status:   req.Status,
	}).Error; err != nil {
		log.Printf("Error updating user ID %s: %v", id, err)
		jsonResponse(w, 500, "Error updating user", nil)
		return
	}

	jsonResponse(w, 0, "User updated successfully!", nil)
}

// API: Lấy danh sách User
func GetListUsers(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// Lấy query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	status := r.URL.Query().Get("status")
	gender := r.URL.Query().Get("gender")
	// Gán vào request model
	req := request.GetListUsersRequest{
		Page:   page,
		Limit:  limit,
		Status: status,
		Gender: gender,
	}

	// Validate dữ liệu
	if err := request.ValidateRequest(&req); err != nil {
		jsonResponse(w, 400, err.Error(), nil)
		return
	}

	// Truy vấn database
	var total int64
	query := db.Model(&models.Users{})
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.Status != "" {
		query = query.Where("gender = ?", req.Gender)
	}
	query.Count(&total)

	var users []models.Users
	offset := (req.Page - 1) * req.Limit
	query.Offset(offset).Limit(req.Limit).Find(&users)

	// Trả về kết quả
	meta := models.Meta{Page: req.Page, Limit: req.Limit, Total: int(total)}
	jsonResponseMeta(w, 0, "User list retrieved successfully", users, meta)
}

// API: Lấy thông tin người dùng theo ID
func GetUserByID(w http.ResponseWriter, idStr string, db *gorm.DB) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		jsonResponse(w, http.StatusBadRequest, "Invalid ID format", nil)
		return
	}

	req := request.GetUserByIDRequest{ID: id}
	if err := request.ValidateRequest(req); err != nil {
		log.Printf("Validation failed: %v", err)
		jsonResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	var user models.Users
	if err := db.First(&user, "id = ?", id).Error; err != nil {
		log.Printf("User not found with ID: %d", id)
		jsonResponse(w, http.StatusNotFound, "User not found", nil)
		return
	}

	jsonResponse(w, 0, "User retrieved successfully", user)
}

// API: Xóa người dùng theo ID
func DeleteUser(w http.ResponseWriter, id string, db *gorm.DB) {
	// Chuyển ID thành số nguyên
	idInt, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("Invalid user ID: %s", id)
		jsonResponse(w, 400, "Invalid ID", nil)
		return
	}

	// Xóa dữ liệu
	if err := db.Delete(&models.Users{}, idInt).Error; err != nil {
		log.Printf("Error deleting user with ID %d: %v", idInt, err)
		jsonResponse(w, 500, "Error deleting user", nil)
		return
	}

	jsonResponse(w, 0, "User deleted successfully", nil)
}

// Hàm JSON response có metadata
func jsonResponseMeta(w http.ResponseWriter, errorCode int, message string, data interface{}, meta models.Meta) {
	w.Header().Set("Content-Type", "application/json")
	response := models.Response{
		ErrorCode: errorCode,
		Message:   message,
		Data:      data,
		Meta:      &meta,
	}
	json.NewEncoder(w).Encode(response)
}

// Hàm JSON response
func jsonResponse(w http.ResponseWriter, errorCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	response := models.Response{
		ErrorCode: errorCode,
		Message:   message,
		Data:      data,
	}
	json.NewEncoder(w).Encode(response)
}
