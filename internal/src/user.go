package src

import (
	"errors"
	"fmt"
	"my_app/internal/db"
	"my_app/internal/logger"
	"my_app/internal/models"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	models.UserModel
}

// 保存数据
func (u *User) Save() error {
	if u == nil {
		return nil
	}
	u.Coin = GetCoin(u.ID)
	return db.DB.Save(u.UserModel).Error
}

// 增加金币
func AddCoin(uid uint, coin int64) int64 {
	key := fmt.Sprintf("coin_%d", uid)
	v, err := db.RedisClient.IncrBy(key, coin)
	if GetCoin(uid) < coin {
		logger.Error(fmt.Sprintf("AddCoin failed not enough coin, uid=%d, coin=%d", uid, coin))
		return 0
	}
	if err != nil {
		logger.Error(fmt.Sprintf("AddCoin failed, uid=%d, coin=%d", uid, coin))
	}
	return v
}

// 获取金币
func GetCoin(uid uint) int64 {
	key := fmt.Sprintf("coin_%d", uid)
	v, err := db.RedisClient.Get(key)
	if err != nil {
		logger.Error(fmt.Sprintf("GetCoin failed, uid=%d", uid))
	}
	coin, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		logger.Error(fmt.Sprintf("GetCoin strconv failed, uid=%d", uid))
	}
	return coin
}

// 通过账号获取玩家数据
func GetUserByName(name, password string) (*User, error) {
	user := models.UserModel{Name: name}
	tx := db.DB.Where(&user).First(&user)
	if tx.RowsAffected == 0 {
		return nil, nil
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("password error")
	}
	return &User{user}, nil
}

// 创建玩家
func CreateUser(name, password string) *User {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(errors.New("error hashing password"))
	}

	user := models.UserModel{
		Name:       name,
		Password:   string(hashedPassword),
		FirstLogin: time.Now(),
		LastLogin:  time.Now(),
		Avatar:     "1",
		RoomID:     0,
	}
	db.DB.Create(&user)
	AddCoin(user.ID, 100)
	return &User{
		UserModel: user,
	}
}
