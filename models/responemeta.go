package models

type ResponseMeta struct {
	ErrorCode int         `json:"error_code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	Meta      Meta        `json:"meta"`
}
