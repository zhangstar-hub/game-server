package models

import (
	"errors"
	"fmt"
	"my_app/internal/db"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID         uint   `gorm:"primary_key"`
	Name       string `gorm:"unique"`
	Password   string
	FirstLogin time.Time
	LastLogin  time.Time
	Coin       uint64
}

// 是否是新用户
func (u *User) IsNewUser() bool {
	return time.Since(u.FirstLogin) <= 7*time.Hour
}

// 保存数据
func (u *User) Save() error {
	if u == nil {
		return nil
	}
	return db.DB.Save(u).Error
}

func GetUserByName(name, password string) (*User, error) {
	user := User{Name: name}
	tx := db.DB.Where(&user).First(&user)
	if tx.RowsAffected == 0 {
		return nil, nil
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("password error")
	}
	return &user, nil
}

func CreateUser(name, password string) *User {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(errors.New("error hashing password"))
	}

	user := &User{
		Name:       name,
		Password:   string(hashedPassword),
		FirstLogin: time.Now(),
		LastLogin:  time.Now(),
	}
	fmt.Printf("user: %#v\n", user)
	db.DB.Create(user)
	return user
}
