package repository

import (
	"gorm.io/gorm"

	"Intent/models"
	"Intent/request"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.Users) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) GetByID(id string) (*models.Users, error) {
	var user models.Users
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(id int, req request.UserUpdateRequest) error {
	return r.db.Model(&models.Users{}).Where("id = ?", id).Updates(req).Error
}

func (r *UserRepository) Delete(id int) error {
	return r.db.Delete(&models.Users{}, id).Error
}
func (r *UserRepository) GetList(req request.GetListUsersRequest) ([]models.Users, int64, error) {
	var users []models.Users
	var total int64 // Khai báo biến total

	query := r.db.Model(&models.Users{})

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.Gender != "" {
		query = query.Where("gender = ?", req.Gender)
	}

	// Đếm tổng số user
	query.Count(&total)

	// Phân trang
	offset := (req.Page - 1) * req.Limit
	err := query.Offset(offset).Limit(req.Limit).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
