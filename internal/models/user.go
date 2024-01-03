package models

import (
	"my_app/internal/db"
	"time"
)

type User struct {
	ID         uint `gorm:"primary_key"`
	Name       string
	password   string
	FirstLogin time.Time
	LastLogin  time.Time
	Coin       uint64
}

func GetUserByName(name, password string) *User {
	user := User{
		Name:     name,
		password: password,
	}
	tx := db.DB.Where(&user).First(&user)
	if tx.RowsAffected == 0 {
		return nil
	}
	return &user
}

func CreateUser(name, password string) *User {
	user := User{
		Name:       name,
		password:   password,
		FirstLogin: time.Now(),
		LastLogin:  time.Now(),
	}
	db.DB.Create(user)
	return &user
}
