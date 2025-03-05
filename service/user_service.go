package service

import (
	"errors"
	"strconv"

	"gorm.io/gorm"

	"Http_Management/models"
	"Http_Management/repository"
	"Http_Management/request"
)

// UserService định nghĩa các chức năng xử lý user
type UserService struct {
	UserRepo *repository.UserRepository
}

// NewUserService khởi tạo service
func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{UserRepo: userRepo}
}

// CreateUser thêm user mới
func (s *UserService) CreateUser(req request.UserCreateRequest) (*models.Users, error) {
	user := models.Users{
		Name:     req.Name,
		Password: req.Password,
		Email:    req.Email,
		Phone:    req.Phone,
		Gender:   req.Gender,
		Status:   req.Status,
	}

	err := s.UserRepo.Create(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser cập nhật thông tin user
func (s *UserService) UpdateUser(id int, req request.UserUpdateRequest) error {
	_, err := s.UserRepo.GetByID(strconv.Itoa(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("User not found")
		}
		return err
	}

	err = s.UserRepo.Update(id, req)
	if err != nil {
		return err
	}

	return nil
}

// GetListUsers lấy danh sách user
func (s *UserService) GetListUsers(req request.GetListUsersRequest) ([]models.Users, models.Meta, error) {
	users, total, err := s.UserRepo.GetList(req)
	if err != nil {
		return nil, models.Meta{}, err
	}

	meta := models.Meta{
		Page:  req.Page,
		Limit: req.Limit,
		Total: total,
	}

	return users, meta, nil
}

// GetUserByID lấy thông tin user theo ID
func (s *UserService) GetUserByID(id int) (*models.Users, error) {
	return s.UserRepo.GetByID(strconv.Itoa(id))
}

// DeleteUser xóa user theo ID
func (s *UserService) DeleteUser(id int) error {
	user, err := s.UserRepo.GetByID(strconv.Itoa(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("User not found")
		}
		return err
	}

	// Xóa user
	err = s.UserRepo.Delete(user.ID)
	if err != nil {
		return err
	}

	return nil
}
