package models

import (
	"time"
)

type UserModel struct {
	ID         uint   `gorm:"primary_key"`
	Name       string `gorm:"unique"`
	Password   string
	FirstLogin time.Time
	LastLogin  time.Time
	Coin       int64
	Avatar     string
	RoomID     uint32
}

// 是否是新用户
func (u *UserModel) IsNewUser() bool {
	return time.Since(u.FirstLogin) <= 7*time.Hour
}
