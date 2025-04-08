package message_queue

import (
	"encoding/json"
	"log"

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

	// Chỉ log lại message
	log.Printf("[Kafka Log] Received action: %s for  user: %s (request_id: %s)", wrapper.Method, data.Email, wrapper.RequestID)
}
