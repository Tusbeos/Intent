package request

type UserCreateRequest struct {
	Name     string `json:"name" validate:"required,min=3"`
	Password string `json:"password" validate:"required,min=6"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required,e164"`
	Gender   string `json:"gender" validate:"required,oneof=male female other"`
	Status   string `json:"status" validate:"required,oneof=active inactive"`
}

type UserUpdateRequest struct {
	ID       int    `json:"id" validate:"required"`
	Name     string `json:"name" validate:"required,min=3"`
	Password string `json:"password" validate:"required,min=6"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required,e164"`
	Gender   string `json:"gender" validate:"required,oneof=male female other"`
	Status   string `json:"status" validate:"required,oneof=active inactive"`
}

func (u *UserCreateRequest) Validate() error {
	return ValidateRequest(u)
}
