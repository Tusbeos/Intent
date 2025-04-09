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
	Method    string                    `json:"method"`
	Path      string                    `json:"path"`
	Request   request.UserCreateRequest `json:"request"`
	RequestID string                    `json:"request_id"`
	Response  string                    `json:"response"`
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

	if wrapper.Method == "" || wrapper.RequestID == "" || wrapper.Request.Email == "" {
		log.Printf("Skipping message: incomplete data. Message: %+v", wrapper)
		return
	}

	log.Printf("Received action: %s for user: %s (request_id: %s)", wrapper.Method, wrapper.Request.Email, wrapper.RequestID)
}
