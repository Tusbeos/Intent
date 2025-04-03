package message_queue

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"gorm.io/gorm"

	"intent/models"
	"intent/repository"
	"intent/request"
)

type Processor struct {
	userRepo *repository.UserRepository
	retries  int
}

func NewProcessor(userRepo *repository.UserRepository, retries int) *Processor {
	return &Processor{
		userRepo: userRepo,
		retries:  retries,
	}
}

func (p *Processor) ProcessMessage(msg []byte) {
	var data request.UserCreateRequest
	if err := json.Unmarshal(msg, &data); err != nil {
		log.Println("[Processor] Failed to unmarshal message:", err)
		return
	}

	// Validate dữ liệu request
	if err := request.ValidateRequest(data); err != nil {
		log.Printf("[Processor] Validation failed: %v", err)
		return
	}
	user := models.Users{
		Name:     data.Name,
		Password: data.Password,
		Email:    data.Email,
		Phone:    data.Phone,
		Gender:   data.Gender,
		Status:   data.Status,
	}

	// Thử lưu user vào database
	for attempt := 1; attempt <= p.retries; attempt++ {
		err := p.userRepo.Create(&user)
		if err == nil {
			log.Println("[Processor] Successfully processed message")
			return
		}

		// Nếu lỗi Duplicate Entry, bỏ qua retry
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			log.Println("[Processor] Skipping duplicate entry")
			return
		}

		log.Printf("[Processor] Attempt %d/%d failed, error: %v", attempt, p.retries, err)
		time.Sleep(2 * time.Second)
	}

	log.Println("[Processor] Failed to process message after retries")
}
