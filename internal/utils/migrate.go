package utils

import (
	"my_app/internal/db"
	"my_app/internal/models"
)

func MirateTable() {
	db.DB.AutoMigrate(&models.User{}, &models.Mission{})
}
