package models

type UserActionMessage struct {
	UserID    int    `json:"user_id"`
	Action    string `json:"action"`
	Timestamp string `json:"timestamp"`
}
