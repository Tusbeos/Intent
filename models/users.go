package models

type Users struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Gender   string `json:"gender"`
	Status   string `json:"status"`
}
