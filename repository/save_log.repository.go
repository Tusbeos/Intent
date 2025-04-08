package repository

import (
	"log"

	"intent/models"
)

func (r *UserRepository) SaveLogAction(logAction models.LogAction) error {
	if err := r.db.Create(&logAction).Error; err != nil {
		log.Println("Failed to save log action:", err)
		return err
	}
	return nil
}
