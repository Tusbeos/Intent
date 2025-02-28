package main

import (
	config "Intent/Config"
	models "Intent/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var cfg *config.Config

// Kết nối Database
func connectDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
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
	http.HandleFunc("PUT /users", func(w http.ResponseWriter, r *http.Request) {
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
	// Khởi động server
	http.ListenAndServe(":8080", nil)
}

// API: Them user
func AddUser(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var user models.Users

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		jsonResponse(w, 400, "Invalid JSON data", nil)
		return
	}
	if err := db.Create(&user).Error; err != nil {
		jsonResponse(w, 500, "Error adding user", nil)
		return
	}

	jsonResponse(w, 0, "User added successfully!", user)
}

// API: Cập nhật người dùng
func UpdateUser(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var user models.Users
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		jsonResponse(w, 400, "Invalid data", nil)
		return
	}

	// Cập nhật dữ liệu
	if err := db.Model(&models.Users{}).Where("id = ?", user.ID).Updates(user).Error; err != nil {
		jsonResponse(w, 500, "Error updating user", nil)
		return
	}

	jsonResponse(w, 0, "User updated successfully!", nil)
}
func GetListUsers(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	// Đếm tổng số user
	var total int64
	db.Model(&models.Users{}).Count(&total)

	// Lấy danh sách user theo phân trang
	var users []models.Users
	offset := (page - 1) * limit
	db.Offset(offset).Limit(limit).Find(&users)

	// Tạo object meta
	meta := models.Meta{
		Page:  page,
		Limit: limit,
		Total: int(total),
	}

	// Trả về JSON response
	jsonResponseMeta(w, 0, "User list retrieved successfully", users, meta)
}

// API: Lấy thông tin người dùng theo ID
func GetUserByID(w http.ResponseWriter, id string, db *gorm.DB) {
	var user models.Users
	if err := db.First(&user, "id = ?", id).Error; err != nil {
		jsonResponse(w, 404, "User not found", nil)
		return
	}
	jsonResponse(w, 0, "User retrieved successfully", user)
}

// API: Xóa người dùng theo ID
func DeleteUser(w http.ResponseWriter, id string, db *gorm.DB) {
	// Chuyển ID thành số nguyên
	idInt, err := strconv.Atoi(id)
	if err != nil {
		jsonResponse(w, 400, "Invalid ID", nil)
		return
	}

	// Xóa dữ liệu
	if err := db.Delete(&models.Users{}, idInt).Error; err != nil {
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
