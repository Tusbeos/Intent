package log_action

import (
	"log"

	"intent/models"
	"intent/repository"
)

func SaveLog(userRepo *repository.UserRepository, logAction models.LogAction) {
	err := userRepo.SaveLogAction(logAction)
	if err != nil {
		log.Println("Failed to save log:", err)
	}
}
