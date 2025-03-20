package models

import "time"

// Struct cho bảng log_actions
type LogAction struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    int    `gorm:"index"`
	Action    string `gorm:"type:varchar(50)"`
	Timestamp time.Time
}
