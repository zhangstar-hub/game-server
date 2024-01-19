package models

import "my_app/internal/db"

func init() {
	db.MirateTable(&User{}, &Mission{}, &LoginBonus{})
}
