package repository

import (
	"Condogenius-notifications/models"

	"gorm.io/gorm"
)

func Save(db *gorm.DB, notification models.Notification) error {
	result := db.Create(&notification)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
