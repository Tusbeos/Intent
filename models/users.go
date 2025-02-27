package models

type Users struct {
	ID       int    `json:"id,omitempty"`
	Name     string `json:"name"`
	Password string `json:"password"`
}
