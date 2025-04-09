package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

// CreateUser
func (r *UserRepository) Create(user *models.Users) error {
	return r.db.Create(user).Error
}

// UpdateUser
func (r *UserRepository) Update(id int, req request.UserUpdateRequest) error {
	return r.db.Debug().Model(&models.Users{}).Where("id = ?", id).Updates(req).Error
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

	val, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var user models.Users
		json.Unmarshal([]byte(val), &user)
		return &user, nil
	}

	var user models.Users
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}

	// Lưu vào cache (1 phút)
	userJSON, _ := json.Marshal(user)
	r.redisClient.Set(ctx, cacheKey, userJSON, 1*time.Minute)

	return &user, nil
}

// CreateBatch tạo nhiều user cùng lúc
func (r *UserRepository) CreateBatch(users []models.Users) error {
	return r.db.Create(&users).Error
}

// FindByEmailOrPhone
func (r *UserRepository) FindByEmailOrPhone(email, phone string) (*models.Users, error) {
	var user models.Users
	if err := r.db.Where("email = ? OR phone = ?", email, phone).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
