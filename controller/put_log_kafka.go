package controller

import (
	"log"

	"intent/request"
)

type MessageWrapper struct {
	Method    string                    `json:"method"`
	Path      string                    `json:"path"`
	Request   request.UserCreateRequest `json:"request"`
	RequestID string                    `json:"request_id"`
	Response  string                    `json:"response"`
}

func (uc *UserController) LogUserActionToKafka(method, path, requestID string, req request.UserCreateRequest) {
	message := MessageWrapper{
		Method:    method,
		Path:      path,
		Request:   req,
		RequestID: requestID,
		Response:  "",
	}
	err := uc.KafkaProducer.Send(message)
	if err != nil {
		log.Printf("Failed to send message with request_id: %s, Error: %v", requestID, err)
	} else {
		log.Printf("Sent message with request_id: %s", requestID)
	}
}
