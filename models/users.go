package models

type Users struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	Name     string `json:"name"`
	Password string `json:"password"`
}
