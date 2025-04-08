package message_queue

import (
	"encoding/json"
	"errors"
	"log"
	"strings"
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
type MessageWrapper struct {
	Method    string                      `json:"method"`
	Path      string                      `json:"path"`
	Request   []request.UserCreateRequest `json:"request"`
	RequestID string                      `json:"request_id"`
	Response  string                      `json:"response"`
}

func NewProcessor(userRepo *repository.UserRepository, retries int) *Processor {
	return &Processor{
		userRepo: userRepo,
		retries:  retries,
	}
}

func (p *Processor) ProcessMessage(msg []byte) {
	var wrapper MessageWrapper
	if err := json.Unmarshal(msg, &wrapper); err != nil {
		log.Printf("Unmarshal failed. Raw message: %s, Error: %v", string(msg), err)
		return
	}

	if len(wrapper.Request) == 0 {
		log.Println("Empty request array in message, skipping.")
		return
	}

	data := wrapper.Request[0]

	user := models.Users{
		Name:     data.Name,
		Password: data.Password,
		Email:    data.Email,
		Phone:    data.Phone,
		Gender:   data.Gender,
		Status:   data.Status,
	}

	for attempt := 1; attempt <= p.retries; attempt++ {
		err := p.userRepo.Create(&user)
		if err == nil {
			log.Println("Successfully processed message")
			return
		}

		if errors.Is(err, gorm.ErrDuplicatedKey) {
			log.Println("Duplicate entry, skipping:", data.Email, data.Phone)
			return
		}

		if strings.Contains(err.Error(), "Data truncated") {
			log.Printf("Invalid data (truncated), skipping. Error: %v, Message: %s", err, string(msg))
			return
		}

		log.Printf("Attempt %d/%d failed for message: %s\nError: %v", attempt, p.retries, string(msg), err)
		time.Sleep(2 * time.Second)
	}

	log.Printf("Failed to process message after retries. Raw message: %s", string(msg))
}
