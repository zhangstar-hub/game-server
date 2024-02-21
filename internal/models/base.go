package models

import "my_app/internal/db"

func MirateTable() {
	db.DB.AutoMigrate(
		&UserModel{}, &MissionModel{}, &LoginBonusModel{},
	)
}
