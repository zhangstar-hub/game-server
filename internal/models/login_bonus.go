package models

import "time"

type LoginBonus struct {
	ID         uint `gorm:"primary_key"`
	SiginTimes int
	LastLogin  time.Time
}

func (p *LoginBonus) InitData() LoginBonus {
	return LoginBonus{
		SiginTimes: 0,
		LastLogin:  time.Now().Add(-24 * time.Hour),
	}
}
