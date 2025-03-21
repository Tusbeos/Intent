package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"intent/models"
	"intent/request"
)

type UserRepository struct {
	db          *gorm.DB
	redisClient *redis.Client
}

func NewUserRepository(db *gorm.DB, redisClient *redis.Client) *UserRepository {
	return &UserRepository{db: db, redisClient: redisClient}
}

// AddUser
func (r *UserRepository) Create(user *models.Users) error {
	return r.db.Create(user).Error
}

// UpdateUser
func (r *UserRepository) Update(id int, req request.UserUpdateRequest) error {
	return r.db.Model(&models.Users{}).Where("id = ?", id).Updates(req).Error
}

// DeleteUser
func (r *UserRepository) Delete(id int) error {
	return r.db.Delete(&models.Users{}, id).Error
}

// GetListUser
func (r *UserRepository) GetList(req request.GetListUsersRequest) ([]models.Users, int64, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("users:status=%s:gender=%s:page=%d:limit=%d", req.Status, req.Gender, req.Page, req.Limit)

	// Kiểm tra cache
	val, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var cachedData struct {
			Users []models.Users
			Total int64
		}
		json.Unmarshal([]byte(val), &cachedData)
		return cachedData.Users, cachedData.Total, nil
	}

	// Nếu không có cache, truy vấn database
	var users []models.Users
	var total int64
	query := r.db.Model(&models.Users{})

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.Gender != "" {
		query = query.Where("gender = ?", req.Gender)
	}

	query.Count(&total)

	offset := (req.Page - 1) * req.Limit
	err = query.Offset(offset).Limit(req.Limit).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	// Lưu vào cache (1 phút)
	cachedData, _ := json.Marshal(struct {
		Users []models.Users
		Total int64
	}{users, total})

	r.redisClient.Set(ctx, cacheKey, cachedData, 1*time.Minute)

	return users, total, nil
}

// GetUserById
func (r *UserRepository) GetByID(id string) (*models.Users, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("user:%s", id)

	// Kiểm tra cache trong Redis
	val, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var user models.Users
		json.Unmarshal([]byte(val), &user)
		return &user, nil
	}

	// Nếu không có cache, truy vấn DB
	var user models.Users
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}

	// Lưu vào cache (1 phút)
	userJSON, _ := json.Marshal(user)
	r.redisClient.Set(ctx, cacheKey, userJSON, 1*time.Minute)

	return &user, nil
}

// Lưu log vào DB
func (r *UserRepository) SaveLogAction(logAction models.LogAction) error {
	if err := r.db.Create(&logAction).Error; err != nil {
		log.Println("Failed to save log action:", err)
		return err
	}
	return nil
}
