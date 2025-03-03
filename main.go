package main

import (
	config "Intent/config"
	models "Intent/models"
	request "Intent/request"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	cfg = config.LoadConfig() // Dùng Viper thay vì hardcode
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

	user := models.Users{Name: req.Name, Password: req.Password}
	if err := db.Create(&user).Error; err != nil {
		log.Printf("Error adding user: %v", err)
		jsonResponse(w, 500, "Error adding user", nil)
		return
	}

	jsonResponse(w, 0, "User added successfully!", user)
}

// API: Cập nhật user
func UpdateUser(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// Lấy ID từ URL
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
	}).Error; err != nil {
		log.Printf("Error updating user ID %s: %v", id, err)
		jsonResponse(w, 500, "Error updating user", nil)
		return
	}

	jsonResponse(w, 0, "User updated successfully!", existingUser)
}

// API: Lấy danh sách User
func GetListUsers(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	req := request.GetListUsersRequest{
		Page:  page,
		Limit: limit,
	}

	// Validate request
	if errMsg := request.ValidateRequest(&req); errMsg != nil {
		log.Printf("Validation error: %v", errMsg)
		jsonResponse(w, 400, errMsg.Error(), nil)
		return
	}

	// Đếm tổng số user
	var total int64
	db.Model(&models.Users{}).Count(&total)

	// Lấy danh sách user theo phân trang
	var users []models.Users
	offset := (req.Page - 1) * req.Limit
	db.Offset(offset).Limit(req.Limit).Find(&users)

	// Tạo object meta
	meta := models.Meta{
		Page:  req.Page,
		Limit: req.Limit,
		Total: int(total),
	}

	// Trả về JSON response
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

	jsonResponse(w, http.StatusOK, "User retrieved successfully", user)
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
