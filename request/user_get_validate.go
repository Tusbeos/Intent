package request

// Request cho API Get List Users
type GetListUsersRequest struct {
	Page   int    `json:"page" validate:"gte=1"`
	Limit  int    `json:"limit" validate:"gte=1,lte=100"`
	Status string `json:"status,omitempty" validate:"omitempty,oneof=active inactive"`
	Gender string `json:"gender,omitempty" validate:"omitempty,oneof=male female"`
}

// Request cho API Get Users By Id
type GetUserByIDRequest struct {
	ID int `json:"id" validate:"required,gt=0"`
}
