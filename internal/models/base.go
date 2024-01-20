package models

import "my_app/internal/db"

func MirateTable() {
	db.MirateTable(&User{}, &Mission{}, &LoginBonus{})
}
