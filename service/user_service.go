package service

import (
	"errors"
	"fmt"
	"strconv"

	"gorm.io/gorm"

	"intent/models"
	"intent/repository"
	"intent/request"
)

// UserService định nghĩa các chức năng xử lý user
type UserService struct {
	UserRepo *repository.UserRepository
}

// NewUserService khởi tạo service
func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{UserRepo: userRepo}
}

// CreateUser
func (s *UserService) CreateUser(req request.UserCreateRequest) (*models.Users, error) {
	existingUser, _ := s.UserRepo.FindByEmailOrPhone(req.Email, req.Phone)
	if existingUser != nil {
		return nil, fmt.Errorf("User with email %s or phone %s already exists", req.Email, req.Phone)
	}

	user := models.Users{
		Name:     req.Name,
		Password: req.Password,
		Email:    req.Email,
		Phone:    req.Phone,
		Gender:   req.Gender,
		Status:   req.Status,
	}

	if err := s.UserRepo.Create(&user); err != nil {
		fmt.Println("Failed to create user:", err)
		return nil, err
	}

	fmt.Println("Successfully created user with ID:", user.ID)
	return &user, nil
}

// UpdateUsers
func (s *UserService) UpdateUser(reqs []request.UserUpdateRequest) error {
	for _, req := range reqs {
		_, err := s.UserRepo.GetByID(strconv.Itoa(req.ID))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("User not found: ID " + strconv.Itoa(req.ID))
			}
			return err
		}

		err = s.UserRepo.Update(req.ID, req)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetListUsers
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

// GetUserByID
func (s *UserService) GetUserByID(id int) (*models.Users, error) {
	return s.UserRepo.GetByID(strconv.Itoa(id))
}

// DeleteUser
func (s *UserService) DeleteUser(id int) error {
	user, err := s.UserRepo.GetByID(strconv.Itoa(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("User not found")
		}
		return err
	}

	err = s.UserRepo.Delete(user.ID)
	if err != nil {
		return err
	}

	return nil
}
