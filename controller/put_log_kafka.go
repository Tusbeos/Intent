package controller

import (
	"log"
)

type MessageWrapper struct {
	Method    string      `json:"method"`
	Path      string      `json:"path"`
	Request   interface{} `json:"request"`
	RequestID string      `json:"request_id"`
	Response  string      `json:"response"`
}

func (uc *UserController) LogUserActionToKafka(method, path, requestID string, req interface{}) {
	message := MessageWrapper{
		Method:    method,
		Path:      path,
		Request:   req,
		RequestID: requestID,
		Response:  "",
	}

	err := uc.KafkaProducer.Send(message)
	if err != nil {
		log.Printf("Failed to send Kafka message | request_id: %s | error: %v\n", requestID, err)
		return
	}

	log.Printf("Kafka message sent successfully. request_id: %s\n", requestID)
}
