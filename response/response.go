package response

import "Http_Management/models"

// SuccessResponse tạo response thành công
func SuccessResponse(code int, message string, data interface{}) models.Response {
	return models.Response{
		ErrorCode: code,
		Message:   message,
		Data:      data,
	}
}

// SuccessResponseWithMeta tạo response thành công với metadata
func SuccessResponseWithMeta(code int, message string, data interface{}, meta models.Meta) models.Response {
	return models.Response{
		ErrorCode: code,
		Message:   message,
		Data:      data,
		Meta:      &meta,
	}
}

// ErrorResponse tạo response lỗi
func ErrorResponse(code int, message string, data interface{}) models.Response {
	return models.Response{
		ErrorCode: code,
		Message:   message,
		Data:      data,
	}
}
