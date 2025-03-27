package repository

import (
	"intent/models"
	"log"
)

// Lưu log vào DB
func (r *UserRepository) SaveLogAction(logAction models.LogAction) error {
	if err := r.db.Create(&logAction).Error; err != nil {
		log.Println("Failed to save log action:", err)
		return err
	}
	return nil
}
