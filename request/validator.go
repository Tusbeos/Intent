package request

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// Custom message lỗi
var customMessages = map[string]string{
	"required": "không được để trống",
	"min":      "phải có ít nhất %s ký tự",
	"gte":      "phải lớn hơn hoặc bằng %s",
	"lte":      "phải nhỏ hơn hoặc bằng %s",
}

// Hàm validate request
func ValidateRequest(r interface{}) error {
	err := validate.Struct(r)
	if err == nil {
		return nil
	}

	var errorMessages []string
	for _, err := range err.(validator.ValidationErrors) {
		fieldName := err.Field()
		tag := err.Tag()
		param := err.Param()

		// Tạo message tùy chỉnh
		if msg, exists := customMessages[tag]; exists {
			if param != "" {
				errorMessages = append(errorMessages, fmt.Sprintf("Trường '%s' %s", fieldName, fmt.Sprintf(msg, param)))
			} else {
				errorMessages = append(errorMessages, fmt.Sprintf("Trường '%s' %s", fieldName, msg))
			}
		} else {
			errorMessages = append(errorMessages, fmt.Sprintf("Trường '%s' không hợp lệ", fieldName))
		}
	}
	return errors.New(strings.Join(errorMessages, "; "))
}
