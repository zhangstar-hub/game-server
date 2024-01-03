package utils

import (
	"my_app/internal/db"
	"my_app/internal/models"
)

func init() {
	db.DB.AutoMigrate(&models.User{})
}
