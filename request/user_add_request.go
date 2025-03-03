package request

// Request struct cho API tạo user
type UserCreateRequest struct {
	Name     string `json:"name" validate:"required,min=3"`
	Password string `json:"password" validate:"required,min=6"`
}

// Request struct cho API cập nhật user
type UserUpdateRequest struct {
	Name     string `json:"name" validate:"omitempty,min=3"`
	Password string `json:"password" validate:"omitempty,min=6"`
}
