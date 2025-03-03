package request

// Request struct cho API tạo user
type UserCreateRequest struct {
	Name     string `json:"name" validate:"required,min=3"`
	Password string `json:"password" validate:"required,min=6"`
	Email    string `json:"email" validate:"required,min=3`
	Phone    string `json:"phone" validate:"required,min=9`
	Gender   string `json:"gender"`
	Status   string `json:"status" validate:"required"`
}

// Request struct cho API cập nhật user
type UserUpdateRequest struct {
	Name     string `json:"name" validate:"required,min=3"`
	Password string `json:"password" validate:"required,min=6"`
	Email    string `json:"email" validate:"required,min=3`
	Phone    string `json:"phone" validate:"required,min=9`
	Gender   string `json:"gender"`
	Status   string `json:"status" validate:"required"`
}
