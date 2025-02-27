package models

type Response struct {
	ErrorCode int         `json:"error_code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
}
